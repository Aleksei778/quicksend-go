package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"quicksend/internal/config"
	"quicksend/internal/gmail"
	tokenmod "quicksend/internal/token"
	usermod "quicksend/internal/user"
	"time"

	"github.com/IBM/sarama"
	"google.golang.org/api/googleapi"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MessagePayload struct {
	Message   string `json:"message"`
	Recipient string `json:"recipient"`
}

type GmailConsumer struct {
	consumer sarama.ConsumerGroup
	userSvc  *usermod.Service
	gmailSvc *gmail.Service
	cfg      *config.Config
}

func NewGmailConsumer(
	consumer sarama.ConsumerGroup,
	userSvc *usermod.Service,
	gmailSvc *gmail.Service,
	cfg *config.Config,
) *GmailConsumer {
	return &GmailConsumer{
		consumer: consumer,
		userSvc:  userSvc,
		gmailSvc: gmailSvc,
		cfg:      cfg,
	}
}

func (g *GmailConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (g *GmailConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (g *GmailConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := g.processMessage(session.Context(), msg); err != nil {
			slog.Error("failed to process message", "err", err)
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

func (g *GmailConsumer) Start(ctx context.Context, topic string) error {
	for {
		if err := g.consumer.Consume(ctx, []string{topic}, g); err != nil {
			return fmt.Errorf("consumer error: %w", err)
		}
		if ctx.Err() != nil {
			return nil
		}
	}
}

func (g *GmailConsumer) processMessage(ctx context.Context, msg *sarama.ConsumerMessage) error {
	var payload MessagePayload
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	userEmail := string(msg.Key)
	user, err := g.userSvc.FindByEmail(userEmail)
	if err != nil {
		slog.Error("consumer: user not found", "email", userEmail)
		return nil
	}

	if err := g.sendWithRetry(ctx, user, payload.Message); err != nil {
		return err
	}

	slog.Info("email sent successfully", "email", userEmail)

	return nil
}

func (g *GmailConsumer) sendWithRetry(ctx context.Context, user *usermod.User, message string) error {
	for attempt := range g.cfg.KafkaMaxRetries {
		_, err := g.gmailSvc.SendEmail(ctx, user, message)
		if err == nil {
			return nil
		}

		if isRetryableError(err) {
			sleep := time.Duration(g.cfg.KafkaBaseBackoff) * time.Duration(math.Pow(2, float64(attempt)))
			select {
			case <-time.After(sleep):
			case <-ctx.Done():
				return ctx.Err()
			}
			continue
		}

		return err
	}

	return fmt.Errorf("sending email failed after %d attempts", g.cfg.KafkaMaxRetries)
}

func isRetryableError(err error) bool {
	if apiErr, ok := errors.AsType[*googleapi.Error](err); ok {
		return apiErr.Code == 429 || apiErr.Code == 500 || apiErr.Code == 505
	}

	return false
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, _ := gorm.Open(postgres.Open(cfg.DSN()))

	saramaCfg := sarama.NewConfig()
	consumer, _ := sarama.NewConsumerGroup(
		[]string{cfg.KafkaBootstrapServers},
		"gmail-consumer-group",
		saramaCfg,
	)
	defer func(consumer sarama.ConsumerGroup) {
		err := consumer.Close()
		if err != nil {
			panic(err)
		}
	}(consumer)

	tokenRepo := tokenmod.NewRepository(db)
	tokenSvc := tokenmod.NewService(db, tokenRepo, cfg)

	gmailSvc := gmail.NewService(tokenSvc, cfg)

	userRepo := usermod.NewRepository(db)
	userSvc := usermod.NewService(db, userRepo)

	gmailConsumer := NewGmailConsumer(consumer, userSvc, gmailSvc, cfg)

	ctx := context.Background()
	if err := gmailConsumer.Start(ctx, cfg.KafkaTopic); err != nil {
		slog.Error("consumer: stopped", "err", err)
	}
}

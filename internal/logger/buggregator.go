package logger

import (
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/goccy/go-json"
)

var levelMap = map[slog.Level]struct {
	level int
	name  string
}{
	slog.LevelDebug: {level: 100, name: "DEBUG"},
	slog.LevelInfo:  {level: 200, name: "INFO"},
	slog.LevelWarn:  {level: 300, name: "WARN"},
	slog.LevelError: {level: 400, name: "ERROR"},
}

type BuggregatorHandler struct {
	conn    net.Conn
	channel string
	opts    slog.HandlerOptions
	attrs   []slog.Attr
}

func NewBuggregatorHandler(conn net.Conn, channel string, opts slog.HandlerOptions) *BuggregatorHandler {
	return &BuggregatorHandler{conn: conn, channel: channel, opts: opts}
}

func (h *BuggregatorHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *BuggregatorHandler) Handle(_ context.Context, r slog.Record) error {
	lvl, ok := levelMap[r.Level]
	if !ok {
		lvl = levelMap[slog.LevelInfo]
	}

	ctx := map[string]any{}
	for _, a := range h.attrs {
		ctx[a.Key] = a.Value.Any()
	}
	r.Attrs(func(a slog.Attr) bool {
		ctx[a.Key] = a.Value.Any()
		return true
	})

	payload := map[string]any{
		"message":    r.Message,
		"context":    ctx,
		"level":      lvl.level,
		"level_name": lvl.name,
		"channel":    h.channel,
		"datetime":   r.Time.Format(time.RFC3339),
		"extra":      map[string]any{},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = h.conn.Write(append(data, '\n'))
	return err
}

func (h *BuggregatorHandler) Close() error {
	return h.conn.Close()
}

func (h *BuggregatorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)
	return &BuggregatorHandler{conn: h.conn, channel: h.channel, opts: h.opts, attrs: newAttrs}
}

func (h *BuggregatorHandler) WithGroup(name string) slog.Handler { return h }

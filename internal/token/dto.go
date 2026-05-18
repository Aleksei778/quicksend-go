package token

import (
	"quicksend/internal/user"
	"time"
)

type FindOrCreate struct {
	User    user.User
	Access  string
	Refresh string
	Expiry  time.Time
}

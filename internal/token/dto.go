package token

import (
	"quicksend/internal/user"
	"time"
)

type FindOrCreate struct {
	user    user.User
	access  string
	refresh string
	expiry  time.Time
}

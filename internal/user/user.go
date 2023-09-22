package user

import (
	"log/slog"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        string    `json:"id,omitempty"`
	Email     string    `json:"email" validate:"required,email"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password,omitempty" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("password hash error", slog.Any("cause", err))
		return httputil.ErrInternalServer("Failed to hash password")
	}
	u.Password = string(hash)
	return nil
}

func (u *User) AssignULID() {
	u.Id = ulid.Make().String()
}

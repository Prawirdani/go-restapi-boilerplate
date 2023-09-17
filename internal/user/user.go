package user

import "time"

type User struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Username  string    `json:"username" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

package auth

import "golang.org/x/crypto/bcrypt"

type LoginRequest struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required"`
}

func (lr *LoginRequest) IsPasswordMatch(storedPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPass), []byte(lr.Password))
	return err == nil
}



package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/prawirdani/go-restapi-boilerplate/internal/user"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/utils"
)

type AuthService interface {
	Login(ctx context.Context, reqLogin LoginRequest) (string, error)
}

type AuthServiceImpl struct {
	userRepository user.UserRepository
	ctxTimeout     time.Duration
}

func NewAuthService(ur user.UserRepository) AuthService {
	return &AuthServiceImpl{userRepository: ur, ctxTimeout: 5 * time.Second}
}

// Login implements UserService.
func (us *AuthServiceImpl) Login(ctx context.Context, reqLogin LoginRequest) (string, error) {
	if err := utils.ValidateRequest(reqLogin); err != nil {
		slog.Error("User.service.req_validator", "cause", err)
		return "", err
	}
	ctxWT, cancel := context.WithTimeout(ctx, us.ctxTimeout)
	defer cancel()

	user, err := us.userRepository.GetUserWithPassword(ctxWT, reqLogin.Email)

	if user == nil || !reqLogin.IsPasswordMatch(user.Password) || err != nil {
		return "", httputil.ErrUnauthorized("check your credentials")
	}

	return "token", nil
}

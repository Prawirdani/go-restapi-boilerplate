package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/prawirdani/go-restapi-boilerplate/internal/user"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/jwt"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/utils"
)

type AuthService interface {
	Login(ctx context.Context, reqLogin LoginRequest) (*jwt.TokenPairs, error)
	RefreshToken(ctx context.Context, userId string) (string, error)
}

type AuthServiceImpl struct {
	userRepository user.UserRepository
	ctxTimeout     time.Duration
}

func NewAuthService(ur user.UserRepository) AuthService {
	return &AuthServiceImpl{userRepository: ur, ctxTimeout: 5 * time.Second}
}

// Login implements UserService.
func (as *AuthServiceImpl) Login(ctx context.Context, reqLogin LoginRequest) (*jwt.TokenPairs, error) {
	if err := utils.ValidateRequest(reqLogin); err != nil {
		slog.Error("User.service.req_validator", "cause", err)
		return nil, err
	}
	
	ctxWT, cancel := context.WithTimeout(ctx, as.ctxTimeout)
	defer cancel()

	user, err := as.userRepository.GetUserWithPassword(ctxWT, reqLogin.Email)
	
	if user == nil || !reqLogin.IsPasswordMatch(user.Password) || err != nil {
		return nil, httputil.ErrUnauthorized("check your credentials")
	}

	tokenPairs := jwt.SignPairs(user)

	return tokenPairs, nil
}

func (as *AuthServiceImpl) RefreshToken(ctx context.Context, userId string) (string, error) {
	ctxWT, cancel := context.WithTimeout(ctx, as.ctxTimeout)
	defer cancel()

	user, err := as.userRepository.GetUserById(ctxWT, userId)
	if err != nil {
		return "", err
	}
	newAccessToken := jwt.SignAccessToken(user)

	return newAccessToken, nil
}

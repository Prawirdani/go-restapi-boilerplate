package user

import (
	"context"
	"log/slog"
	"time"

	"github.com/prawirdani/go-restapi-boilerplate/pkg/utils"
)

type UserService interface {
	FindById(ctx context.Context, id int) (*User, error)
	FindAll(ctx context.Context) ([]User, error)
	Save(ctx context.Context, request User) error
}

type UserServiceImpl struct {
	userRepository UserRepository
	ctxTimeout     time.Duration
}

func NewUserService(ur UserRepository) UserService {
	return &UserServiceImpl{userRepository: ur, ctxTimeout: 5 * time.Second}
}

// FindAll implements UserService.
func (us *UserServiceImpl) FindAll(ctx context.Context) ([]User, error) {
	ctxWT, cancel := context.WithTimeout(ctx, us.ctxTimeout)
	defer cancel()
	return us.userRepository.GetUsers(ctxWT)
}

// FindById implements UserService.
func (us *UserServiceImpl) FindById(ctx context.Context, id int) (*User, error) {
	ctxWT, cancel := context.WithTimeout(ctx, us.ctxTimeout)
	defer cancel()
	return us.userRepository.GetUserById(ctxWT, id)
}

// Save implements UserService.
func (us *UserServiceImpl) Save(ctx context.Context, request User) error {
	ctxWT, cancel := context.WithTimeout(ctx, us.ctxTimeout)
	defer cancel()

	if err := utils.ValidateRequest(request); err != nil {
		slog.Error("User.service.req_validator", "cause", err)
		return err
	}
	request.HashPassword()

	return us.userRepository.CreateUser(ctxWT, request)
}

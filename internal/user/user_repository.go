package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, newUser User) error
	GetUserById(ctx context.Context, id int) (*User, error)
	GetUserWithPassword(ctx context.Context, email string) (*User, error)
}

type UserRepositoryImpl struct {
	postgres *pgxpool.Pool
}

func NewUserRepository(pgConn *pgxpool.Pool) UserRepository {
	return &UserRepositoryImpl{postgres: pgConn}
}

// CreateUser implements UserRepository.
func (ur *UserRepositoryImpl) CreateUser(ctx context.Context, newUser User) error {
	q := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	if _, err := ur.postgres.Exec(ctx, q, newUser.Username, newUser.Email, newUser.Password); err != nil {
		return err
	}
	return nil
}

// GetUserById implements UserRepository.
func (ur *UserRepositoryImpl) GetUserById(ctx context.Context, id int) (*User, error) {
	q := `SELECT id, username, email, created_at from users WHERE id = $1`

	var user User

	if err := ur.postgres.QueryRow(ctx, q, id).Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, httputil.ErrNotFound("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUsers implements UserRepository.
func (ur *UserRepositoryImpl) GetUsers(ctx context.Context) ([]User, error) {
	var users []User

	rows, _ := ur.postgres.Query(ctx, "Select id, username, email, created_at from users LIMIT 100")

	for rows.Next() {
		var each User
		if err := rows.Scan(&each.Id, &each.Username, &each.Email, &each.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, each)
	}

	return users, nil
}

// GetUserWithPassword implements UserRepository.
func (ur *UserRepositoryImpl) GetUserWithPassword(ctx context.Context, email string) (*User, error) {
	var user User
	q := `SELECT id, username, email, password from users WHERE email = $1`

	if err := ur.postgres.QueryRow(ctx, q, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		if err == pgx.ErrNoRows {
			return nil, httputil.ErrNotFound("user not found")
		}
		return nil, err
	}

	return &user, nil
}

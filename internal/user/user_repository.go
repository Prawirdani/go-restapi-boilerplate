package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prawirdani/go-restapi-boilerplate/pkg/httputil"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id int) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
	CreateUser(ctx context.Context, newUser User) error
}

type UserRepositoryImpl struct {
	postgres *pgxpool.Pool
}

func NewUserRepository(pgConn *pgxpool.Pool) UserRepository {
	return &UserRepositoryImpl{postgres: pgConn}
}

// CreateUser implements UserRepository.
func (ur *UserRepositoryImpl) CreateUser(ctx context.Context, newUser User) error {
	q := `INSERT INTO users (username, email) VALUES ($1, $2)`
	if _, err := ur.postgres.Exec(ctx, q, newUser.Username, newUser.Email); err != nil {
		return err
	}
	return nil
}

// GetUserById implements UserRepository.
func (ur *UserRepositoryImpl) GetUserById(ctx context.Context, id int) (*User, error) {
	q := `SELECT * from users WHERE id = $1`

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

	rows, _ := ur.postgres.Query(ctx, "Select * from users")

	for rows.Next() {
		var each User
		if err := rows.Scan(&each.Id, &each.Username, &each.Email, &each.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, each)
	}

	return users, nil
}

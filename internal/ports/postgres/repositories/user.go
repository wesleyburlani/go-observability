package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	db "github.com/wesleyburlani/go-observability/internal/ports/postgres"
	"github.com/wesleyburlani/go-observability/internal/users"
	pkg_errors "github.com/wesleyburlani/go-observability/pkg/errors"
)

type UserRepository struct {
	db *db.Database
}

func NewUserRepository(db *db.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Get(ctx context.Context, id int64) (users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u, err := r.db.Queries.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return users.User{}, fmt.Errorf("could not found user with id %d: %w", id, pkg_errors.ErrNotFound)
		}
		return users.User{}, err
	}

	return r.entityToDTO(u), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u, err := r.db.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return users.User{}, err
	}

	return r.entityToDTO(u), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u, err := r.db.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		return users.User{}, err
	}

	return r.entityToDTO(u), nil
}

func (r *UserRepository) Create(ctx context.Context, u users.User) (users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	nu, err := r.db.Queries.CreateUser(ctx, db.CreateUserParams{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return users.User{}, fmt.Errorf("%w: user already exists", pkg_errors.ErrConflict)
		}
		return users.User{}, err
	}

	return r.entityToDTO(nu), nil
}

func (r *UserRepository) Update(ctx context.Context, u users.User) (users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	pwd := sql.NullString{Valid: u.Password != "", String: u.Password}
	nu, err := r.db.Queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       u.ID,
		Username: sql.NullString{Valid: u.Username != "", String: u.Username},
		Email:    sql.NullString{Valid: u.Email != "", String: u.Email},
		Password: pwd,
	})
	if err != nil {
		return users.User{}, err
	}

	return r.entityToDTO(nu), nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) (users.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u, err := r.db.Queries.DeleteUser(ctx, id)
	if err != nil {
		return users.User{}, err
	}

	return r.entityToDTO(u), nil
}

func (r *UserRepository) entityToDTO(u db.User) users.User {
	return users.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

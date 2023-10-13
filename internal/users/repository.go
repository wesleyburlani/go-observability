package users

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	pkg_errors "github.com/wesleyburlani/go-rest/pkg/errors"

	"github.com/wesleyburlani/go-rest/internal/db"
)

type Repository struct {
	db *db.Database
}

func NewRepository(db *db.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, id int64) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u, err := r.db.Queries.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("could not found user with id %d: %w", id, pkg_errors.ErrNotFound)
		}
		return User{}, err
	}

	return r.entityToDTO(u), nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u, err := r.db.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}

	return r.entityToDTO(u), nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u, err := r.db.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		return User{}, err
	}

	return r.entityToDTO(u), nil
}

func (r *Repository) Create(ctx context.Context, u User) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	encryptedPwd, err := encryptPassword(u.Password)

	if err != nil {
		return User{}, fmt.Errorf("could not encrypt password: %w", err)
	}

	nu, err := r.db.Queries.CreateUser(ctx, db.CreateUserParams{
		Username: u.Username,
		Email:    u.Email,
		Password: encryptedPwd,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return User{}, fmt.Errorf("%w: user already exists", pkg_errors.ErrConflict)
		}
		return User{}, err
	}

	return r.entityToDTO(nu), nil
}

func (r *Repository) Update(ctx context.Context, u User) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pwd := sql.NullString{Valid: u.Password != "", String: u.Password}

	if pwd.Valid {
		encryptedPwd, err := encryptPassword(u.Password)
		if err != nil {
			return User{}, fmt.Errorf("could not encrypt password: %w", err)
		}
		pwd = sql.NullString{Valid: true, String: encryptedPwd}
	}

	nu, err := r.db.Queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       u.ID,
		Username: sql.NullString{Valid: u.Username != "", String: u.Username},
		Email:    sql.NullString{Valid: u.Email != "", String: u.Email},
		Password: pwd,
	})
	if err != nil {
		return User{}, err
	}

	return r.entityToDTO(nu), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u, err := r.db.Queries.DeleteUser(ctx, id)
	if err != nil {
		return User{}, err
	}

	return r.entityToDTO(u), nil
}

func (r *Repository) entityToDTO(u db.User) User {
	return User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

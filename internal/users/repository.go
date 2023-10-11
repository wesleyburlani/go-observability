package users

import (
	"context"
	"database/sql"
	"time"

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

	nu, err := r.db.Queries.CreateUser(ctx, db.CreateUserParams{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		return User{}, err
	}

	return r.entityToDTO(nu), nil
}

func (r *Repository) Update(ctx context.Context, u User) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	nu, err := r.db.Queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       u.ID,
		Username: sql.NullString{Valid: u.Username != "", String: u.Username},
		Email:    sql.NullString{Valid: u.Email != "", String: u.Email},
		Password: sql.NullString{Valid: u.Password != "", String: u.Password},
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

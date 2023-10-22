package users

import (
	"context"

	pkg_errors "github.com/wesleyburlani/go-observability/pkg/errors"
	"github.com/wesleyburlani/go-observability/pkg/logger"
)

type Service struct {
	repo   *Repository
	logger *logger.Logger
}

func NewService(repo *Repository, logger *logger.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) Get(ctx context.Context, id int64) (User, error) {
	s.logger.With("user", id).Debug(ctx, "getting user")
	return s.repo.Get(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (User, error) {
	s.logger.With("email", email).Debug(ctx, "getting user")
	return s.repo.GetByEmail(ctx, email)
}

func (s *Service) GetByUsername(ctx context.Context, username string) (User, error) {
	s.logger.With("username", username).Debug(ctx, "getting user")
	return s.repo.GetByUsername(ctx, username)
}

func (s *Service) Create(ctx context.Context, u User) (User, error) {
	s.logger.With("user", u).Debug(ctx, "creating user")
	return s.repo.Create(ctx, u)
}

func (s *Service) Update(ctx context.Context, u User) (User, error) {
	s.logger.With("user", u).Debug(ctx, "updating user")
	return s.repo.Update(ctx, u)
}

func (s *Service) Delete(ctx context.Context, id int64) (User, error) {
	s.logger.With("user", id).Debug(ctx, "deleting user")
	return s.repo.Delete(ctx, id)
}

func (s *Service) Login(ctx context.Context, username, password string) error {
	s.logger.With("username", username).Debug(ctx, "logging in user")
	u, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return pkg_errors.ErrUnauthorized
	}

	err = compareHashAndPassword(u.Password, password)
	if err != nil {
		return pkg_errors.ErrUnauthorized
	}
	return nil
}

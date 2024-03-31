package users

import (
	"context"
	"fmt"
	"sync"

	pkg_errors "github.com/wesleyburlani/go-observability/pkg/errors"
	"github.com/wesleyburlani/go-observability/pkg/logger"
)

type Service struct {
	observers []UserEventsObserver
	repo      Repository
	logger    *logger.Logger
}

func NewService(repo Repository, logger *logger.Logger, observers []UserEventsObserver) *Service {
	return &Service{repo: repo, logger: logger, observers: observers}
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
	encryptedPwd, err := encryptPassword(u.Password)
	if err != nil {
		return User{}, fmt.Errorf("could not encrypt password: %w", err)
	}
	u.Password = encryptedPwd
	createdUser, err := s.repo.Create(ctx, u)
	if err != nil {
		return User{}, err
	}
	var waitGroup sync.WaitGroup
	for _, observer := range s.observers {
		waitGroup.Add(1)
		go func(observer UserEventsObserver) {
			observer.OnUserCreated(ctx, createdUser)
			waitGroup.Done()
		}(observer)
	}
	waitGroup.Wait()
	return createdUser, nil
}

func (s *Service) Update(ctx context.Context, u User) (User, error) {
	s.logger.With("user", u).Debug(ctx, "updating user")
	encryptedPwd, err := encryptPassword(u.Password)
	if err != nil {
		return User{}, fmt.Errorf("could not encrypt password: %w", err)
	}
	u.Password = encryptedPwd
	updatedUser, err := s.repo.Update(ctx, u)
	if err != nil {
		return User{}, err
	}
	var waitGroup sync.WaitGroup
	for _, observer := range s.observers {
		waitGroup.Add(1)
		go func(observer UserEventsObserver) {
			observer.OnUserUpdated(ctx, updatedUser)
			waitGroup.Done()
		}(observer)
	}
	waitGroup.Wait()
	return updatedUser, nil
}

func (s *Service) Delete(ctx context.Context, id int64) (User, error) {
	s.logger.With("user", id).Debug(ctx, "deleting user")
	deletedUser, err := s.repo.Delete(ctx, id)
	if err != nil {
		return User{}, err
	}
	var waitGroup sync.WaitGroup
	for _, observer := range s.observers {
		waitGroup.Add(1)
		go func(observer UserEventsObserver) {
			observer.OnUserDeleted(ctx, deletedUser)
			waitGroup.Done()
		}(observer)
	}
	waitGroup.Wait()
	return deletedUser, nil
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

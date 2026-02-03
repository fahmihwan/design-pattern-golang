package service

import (
	"best-pattern/internal/model"
	"best-pattern/internal/repository"
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var _ UserServiceInterface = &UserService{}

type UserServiceInterface interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
	Login(ctx context.Context, email, password string) (*model.User, error)
}

type UserService struct {
	repo repository.Repository
}

func NewUserService(repo repository.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(ctx context.Context, user *model.User) (*model.User, error) {

	existingUser, err := s.repo.User.FindByEmail(ctx, user.Email)
	if err == nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %q already exists", user.Email)
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if user.Password == nil {
		return nil, fmt.Errorf("password is required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	hashedStr := string(hash)
	user.Password = &hashedStr

	if err := s.repo.User.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, error) {

	user, err := s.repo.User.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user : %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("Invalid email or password")
	}

	user.Password = nil
	return user, nil

}

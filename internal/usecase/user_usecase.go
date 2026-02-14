package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

type UserUsecase struct {
	userRepo port.UserRepository
}

func NewUserUsecase(userRepo port.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (u *UserUsecase) CreateUser(ctx context.Context, name string) (*domain.User, error) {
	id := uuid.New().String()
	user, err := domain.NewUser(id, name)
	if err != nil {
		return nil, err
	}
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

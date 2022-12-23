package usecases

import (
	"errors"

	"github.com/tkanzakic/cellar/helpers"
	"github.com/tkanzakic/cellar/internal/core/domain"
	"github.com/tkanzakic/cellar/internal/core/ports"
)

type authUseCase struct {
	repo ports.UserRepository
}

func NewAuthUseCase(userRepo ports.UserRepository) ports.AuthUseCase {
	return &authUseCase{
		repo: userRepo,
	}
}

func (u *authUseCase) SignUp(email, name, password string) (*domain.User, error) {
	id := helpers.Uuid()
	user := domain.NewUserHashingPassword(id, email, name, password)
	user, err := u.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUseCase) SignIn(email, password string) (*domain.User, error) {
	user, err := u.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("Invalid email address")
	}
	if !user.VerifyPassword(password) {
		return nil, errors.New("Invalid password")
	}

	return user, nil
}

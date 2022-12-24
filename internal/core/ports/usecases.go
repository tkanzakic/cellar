package ports

import "github.com/tkanzakic/cellar/internal/core/domain"

type AuthUseCase interface {
	SignUp(family, email, name, password string) (*domain.User, error)
	SignIn(family, email, password string) (*domain.User, error)
}

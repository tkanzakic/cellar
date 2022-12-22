package ports

import "github.com/tkanzakic/cellar/internal/core/domain"

type AuthUseCase interface {
	SignUp(email, name, password string) (*domain.User, error)
	SignIn(email, password string) (*domain.User, error)
}

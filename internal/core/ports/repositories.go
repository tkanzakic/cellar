package ports

import "github.com/tkanzakic/cellar/internal/core/domain"

type UserRepository interface {
	GetByEmail(family, email string) (*domain.User, error)
	Create(user *domain.User) (*domain.User, error)
}

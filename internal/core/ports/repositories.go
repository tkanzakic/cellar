package ports

import "github.com/tkanzakic/cellar/internal/core/domain"

type UserRepository interface {
	Get(id string) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Create(user *domain.User) (*domain.User, error)
}

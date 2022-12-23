package repositories

import (
	"encoding/json"
	"errors"

	"github.com/tkanzakic/cellar/internal/core/domain"
)

type userInMemoryRepository struct {
	kvs map[string][]byte
}

func NewInMemoryUserRepository() *userInMemoryRepository {
	return &userInMemoryRepository{kvs: map[string][]byte{}}
}

func (r *userInMemoryRepository) Get(id string) (*domain.User, error) {
	if value, ok := r.kvs[id]; ok {
		user, err := unmarshal(value)
		return &user, err
	}
	return nil, errors.New("User does not exists")
}

func (r *userInMemoryRepository) GetByEmail(email string) (*domain.User, error) {
	for _, value := range r.kvs {
		user, err := unmarshal(value)
		if err == nil && user.Email == email {
			return &user, nil
		}
	}
	return nil, errors.New("User does not exists")
}

func (r *userInMemoryRepository) Create(user *domain.User) (*domain.User, error) {
	value, err := json.Marshal(user)
	if err != nil {
		return nil, errors.New("Error marshalling user")
	}
	r.kvs[user.ID] = value
	return user, nil
}

func unmarshal(value []byte) (domain.User, error) {
	user := domain.User{}
	err := json.Unmarshal(value, &user)
	if err != nil {
		return domain.User{}, errors.New("Could not unmarshal user")
	}
	return user, nil
}

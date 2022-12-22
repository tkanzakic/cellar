package domain

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewUser(id, email, name, password string) *User {
	return &User{
		ID:       id,
		Email:    email,
		Name:     name,
		Password: password,
	}
}

func NewUserHashingPassword(id, email, name, password string) *User {
	pass, err := hashPassword(password)
	if err != nil {
		panic("Cannot hash password")
	}

	return &User{
		ID:       id,
		Email:    email,
		Name:     name,
		Password: pass,
	}
}

func (u *User) HashedPassword() (string, error) {
	bytes, err := hashPassword(u.Password)
	return string(bytes), err
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

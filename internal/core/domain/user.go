package domain

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Family   string
	Email    string
	Name     string
	Password string
}

func NewUser(family, email, name, password string) *User {
	return &User{
		Family:   family,
		Email:    email,
		Name:     name,
		Password: password,
	}
}

func NewUserHashingPassword(family, email, name, password string) *User {
	pass, err := hashPassword(password)
	if err != nil {
		panic("Cannot hash password")
	}

	return NewUser(family, email, name, pass)
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

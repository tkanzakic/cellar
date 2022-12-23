package usecases

import (
	"testing"

	"github.com/tkanzakic/cellar/internal/core/domain"
	"github.com/tkanzakic/cellar/internal/core/ports"
	"github.com/tkanzakic/cellar/internal/repositories"
)

var (
	email    = "email@server.com"
	name     = "User Full Name"
	password = "Pas2sw0rd"
)

// SignUp
func TestShouldCreateUser(t *testing.T) {
	userRepository, sut := givenSut()

	_, err := sut.SignUp(email, name, password)
	if err != nil {
		t.Fatal("Sign up failed")
	}

	_, err = userRepository.GetByEmail(email)
	if err != nil {
		t.Fatal("User not created")
	}
}

func TestShouldNotCreateUserIfEmailAlreadyExists(t *testing.T) {
	_, sut := givenSut()
	givenUserCreated(sut)

	_, err := sut.SignUp(email, "Other name", "Other password")
	if err == nil {
		t.Fatal("Sign up succeed for duplicated email")
	}
}

// SignIn
func TestShouldSignInUser(t *testing.T) {
	_, sut := givenSut()
	givenUserCreated(sut)

	_, err := sut.SignIn(email, password)

	if err != nil {
		t.Fatal("Sign in failed")
	}
}

func TestShouldReturnErrorIfUserDoesNotExists(t *testing.T) {
	_, sut := givenSut()

	user, err := sut.SignIn(email, password)

	if user != nil || err == nil {
		t.Fatal("Sign in succeeded when user does not exists")
	}
}

func TestShouldReturnIfInvalidPassword(t *testing.T) {
	_, sut := givenSut()
	givenUserCreated(sut)

	user, err := sut.SignIn(email, "Invalid password")

	if user != nil || err == nil {
		t.Fatal("Sign in succeeded with invalid password")
	}
}

// Utility functions
func givenSut() (ports.UserRepository, ports.AuthUseCase) {
	userRepository := repositories.NewInMemoryUserRepository()
	return userRepository, NewAuthUseCase(userRepository)
}

func givenUserCreated(useCase ports.AuthUseCase) *domain.User {
	user, err := useCase.SignUp(email, name, password)

	if err != nil {
		panic("Could not create user")
	}
	return user
}

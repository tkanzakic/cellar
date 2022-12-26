package main

import (
	"net/http"

	"github.com/tkanzakic/cellar/internal/core/usecases"
	"github.com/tkanzakic/cellar/internal/handlers/signin"
	"github.com/tkanzakic/cellar/internal/repositories"
)

var signInHandler = signin.NewHTTHandler(usecases.NewAuthUseCase(repositories.NewDynamoDBUserRepository()))

func main() {
	http.HandleFunc("/signin", signInHandler.SignIn)
	http.ListenAndServe(":8080", nil)
}

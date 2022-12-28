package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tkanzakic/cellar/internal/core/usecases"
	"github.com/tkanzakic/cellar/internal/handlers/signin"
	"github.com/tkanzakic/cellar/internal/repositories"
)

var signInHandler = signin.NewHTTHandler(usecases.NewAuthUseCase(repositories.NewDynamoDBUserRepository()))

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Print("PORT not specified, using default one")
		port = "5000"
	} else {
		log.Print("Using custom port: " + port)
	}

	http.HandleFunc("/signin", signInHandler.SignIn)

	http.ListenAndServe(":"+port, nil)
}

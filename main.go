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
	logFile, _ := os.Create("/var/log/golang/golang-server.log")
	log.SetOutput(logFile)
	defer logFile.Close()

	port := os.Getenv("PORT")
	if port == "" {
		log.Default().Print("PORT not specified, using default one")
		port = "5000"
	}

	http.HandleFunc("/signin", signInHandler.SignIn)

	http.ListenAndServe(":"+port, nil)
}

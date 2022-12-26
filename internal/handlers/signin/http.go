package signin

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/tkanzakic/cellar/internal/core/domain"
	"github.com/tkanzakic/cellar/internal/core/ports"
)

var logger = log.Default()

type HTTPHandler struct {
	useCase ports.AuthUseCase
}

type signInRequest struct {
	Family   string `json:"family"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *signInRequest) isValid() bool {
	return r.Email != "" && r.Password != "" && r.Family != ""
}

type signInResponse struct {
	Family string `json:"family"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type signInErrorResponse struct {
	Message string `json:"message"`
}

func NewHTTHandler(useCase ports.AuthUseCase) *HTTPHandler {
	return &HTTPHandler{
		useCase: useCase,
	}
}

func (h *HTTPHandler) SignIn(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var reqBody signInRequest
	err := json.NewDecoder(request.Body).Decode(&reqBody)
	responseEncoder := json.NewEncoder(writer)
	if err != nil || !reqBody.isValid() {
		logger.Printf("Invalid request received")
		writer.WriteHeader(http.StatusBadRequest)
		responseEncoder.Encode(signInErrorResponse{
			Message: "Invalid request",
		})
		return
	}
	user, err := h.useCase.SignIn(reqBody.Family, reqBody.Email, reqBody.Password)
	if err != nil {
		logger.Printf("Invalid credentials request received")
		writer.WriteHeader(http.StatusForbidden)
		responseEncoder.Encode(signInErrorResponse{
			Message: "Invalid credentials",
		})
		return
	}
	token, err := generateJWT(user)
	if err != nil {
		logger.Printf("An error ocurred while generating JWT token.\n%v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		responseEncoder.Encode(signInErrorResponse{
			Message: "An unexpected error ocurred, please try again later",
		})
		return
	}
	logger.Println("🚀 SingIn succeeded")
	writer.Header().Add("X-Jwt-Token", token)
	writer.WriteHeader(http.StatusOK)
	responseEncoder.Encode(signInResponse{
		Family: user.Family,
		Email:  user.Email,
		Name:   user.Name,
	})
}

func generateJWT(user *domain.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 43200) // 30 days
	claims["family"] = user.Family
	claims["email"] = user.Email
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

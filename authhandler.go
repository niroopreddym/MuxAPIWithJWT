package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//AuthHandler handles all the operations related to the JWT Tokens
type AuthHandler struct {
}

//NewAuthHandler returns a new instance of the Auth Handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

//GenerateToken generates a new JWT Token
func (handler *AuthHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	tokenString, err := GenerateJWT()
	if err != nil {
		fmt.Println("error occured while generating the token string")
	}

	fmt.Fprintf(w, tokenString)
}

var mySigningKey = []byte("ultimateStarAjith")

//GenerateJWT generates JWT Token
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["user"] = "Niroop"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return tokenString, nil
}

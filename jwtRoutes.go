package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// User validation
func validateLogin(ctx *gin.Context) {
	// TODO: Do some account validation.
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	fmt.Printf("Username: %v | Password: %v\n", username, password)

	// If the account is authenticated, then generate the token.
	tokenString, err := generateToken(username)
	if err != nil {
		log.Fatal(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func generateToken(username string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	secretKey := os.Getenv("SECRET_KEY")
	secretKeyByte := []byte(secretKey)

	// Create a new JWT token using this signing method.
	token := jwt.New(jwt.SigningMethodHS256)

	// Map out the claims.
	claims := token.Claims.(jwt.MapClaims)
	// claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["aud"] = username

	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Token validation middleware
func validateToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("Authorization")
		if authorizationHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "provide your token using the Authorization header and Bearer schema."})
			ctx.Abort()
			return
		}

		// Split the schema and the value.
		authorizationValue := strings.Split(authorizationHeader, " ")
		tokenString := authorizationValue[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			_, ok := t.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return "", fmt.Errorf("error when parsing")
			}

			secretKey := []byte(os.Getenv("SECRET_KEY"))
			return secretKey, nil
		})

		if err != nil {
			panic(err)
		}

		if token.Valid {
			ctx.Next()
		} else {
			ctx.Abort()
		}
	}
}

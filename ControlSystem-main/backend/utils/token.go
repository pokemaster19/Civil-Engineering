package utils

import (
	"ControlSystem/models"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(user models.User) (string, error) {

	tokenLifespan, err := strconv.Atoi(strings.TrimSpace(os.Getenv("TOKEN_MINUTE_LIFESPAN")))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

func GenerateRefreshToken(user models.User) (string, error) {
	tokenLifespan, err := strconv.Atoi(strings.TrimSpace(os.Getenv("REFRESH_TOKEN_HOUR_LIFESPAN")))
	if err != nil {
		tokenLifespan = 72
		fmt.Println("Warning: REFRESH_TOKEN_HOUR_LIFESPAN not set or invalid. Using default 72 hours.")
	}

	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
}

func ValidateToken(c *gin.Context) error {
	token, err := GetToken(c)

	if err != nil {
		return err
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}

	return errors.New("invalid token provided")
}

func GetToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("API_SECRET")), nil
	})
	return token, err
}

func getTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func ExtractUserIDFromToken(token *jwt.Token) (uint, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token claims")
	}

	idFloat, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("user ID not found or invalid format in token")
	}
	return uint(idFloat), nil
}

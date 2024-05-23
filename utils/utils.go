package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func IsOverlapping(startTime, endTime time.Time) bool {
	// logger.Printf("FUNCTION - isOverlapping, params : %s, %s", startTime, endTime)
	fmt.Println(startTime, endTime)
	return false
}

// VerifyToken verifies a token JWT validate
func VerifyToken(tokenString string) (jwt.Claims, error) {
	// Parse the token
	var secretKey = []byte(os.Getenv("SECRET"))
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["user_id"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

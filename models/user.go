package models

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name  string `json:"name"`
	Email string `gorm:"unique_index" json:"email"`
	ID    string `gorm:"primary_key" json:"id"`
	Hash  string `json:"-"`
}

type JWTToken struct {
	Token string `json:"token"`
}

func (u User) HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes)
}

func (u User) CheckPassword(password string) bool {
	fmt.Printf("User Hash : %s", u.Hash)
	err := bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(password))
	return err == nil
}

func (u User) GenerateJWT() (JWTToken, error) {
	signingKey := []byte(os.Getenv("SECRET"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 1 * 1).Unix(),
		"user_id": u.ID,
		"name":    u.Name,
		"email":   u.Email,
	})
	tokenString, err := token.SignedString(signingKey)
	return JWTToken{tokenString}, err
}

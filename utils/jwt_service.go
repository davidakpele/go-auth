package utils

import (
	"time"
	"os"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(userID uint, email string, role string) (string, error) {
    claims := jwt.MapClaims{
        "id":    userID,
        "email": email,
        "roles": []string{role}, 
        "iso":   time.Now(),              
        "exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }
    return signedToken, nil
}
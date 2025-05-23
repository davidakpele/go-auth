package security

import (
	"api-service/models"
	"api-service/repositories"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// SecurityFilterChain struct
type SecurityFilterChain struct {
	JwtKey         string
	UserRepository *repositories.UserRepository
}

func (s *SecurityFilterChain) IsValidToken(c *gin.Context) (*models.User, error) {
    authHeader := c.GetHeader("Authorization")

    // Ensure Authorization header exists
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  http.StatusUnauthorized,
            "title":   "Authentication Error",
            "details": "Missing Authorization Header",
        })
        return nil, errors.New("missing authorization header")
    }

    // Check Bearer token format
    if !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  http.StatusUnauthorized,
            "title":   "Authentication Error",
            "details": "Invalid Token Format: Bearer token required",
        })
        return nil, errors.New("invalid token format")
    }
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")

    // Decode and validate the token
    user, err := s.VerifyToken(tokenString)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "status":  http.StatusUnauthorized,
            "title":   "Authentication Error",
            "details": "Invalid Token: " + err.Error(),
        })
        return nil, err
    }
    return user, nil
}

func (s *SecurityFilterChain) VerifyToken(tokenString string) (*models.User, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Ensure signing method is HMAC
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(s.JwtKey), nil
    })
    if err != nil {
        return nil, errors.New("invalid token")
    }

    // Extract claims and validate payload
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userEmail, okEmail := claims["sub"].(string)
        userIDFloat, okID := claims["group"].(map[string]interface{})["id"].(float64)

        if !okEmail || !okID {
            return nil, errors.New("invalid token claims structure")
        }

        userID := int(userIDFloat)

        // Verify user in the database
        user, err := s.UserRepository.VerifyUserAccount(userEmail, userID)
        if err != nil || user == nil {
            return nil, errors.New("unauthorized user")
        }

        return user, nil
    }

    return nil, errors.New("invalid token payload")
}

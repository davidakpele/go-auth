package middleware

import (
	"api-service/repositories"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

// Update the middleware to accept userRepo as a parameter
func AuthenticationMiddleware(userRepo *repositories.UserRepository) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error":   "Access Denied",
                "status": "error", 
                "title": "Authentication Error", 
                "message": "Authorization Access",
                "details": "Something went wrong with authentication to your SkyBase library.", 
                "code": "generic_authentication_error",
            })
            c.Abort()
            return
        }

        // Extract token from "Bearer <token>" format
        parts := strings.Fields(authHeader)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid token format"})
            c.Abort()
            return
        }
        tokenString := parts[1]

        // Parse and validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Ensure signing method matches
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            log.Printf("Token validation error: %v", err)
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid or expired token"})
            c.Abort()
            return
        }

        // Extract claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid token claims"})
            c.Abort()
            return
        }

        // Get user ID from claims
        idFloat, ok := claims["id"].(float64)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{
                "status": "error",
                "message": "User ID must be a number",
            })
            c.Abort()
            return
        }
        userID := uint(idFloat)
        // Fetch user from repository
        user, err := userRepo.GetUserByID(userID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "status": "error",
                "message": "User not found",
            })
            c.Abort()
            return
        }

        // Set user information in context
        c.Set("user", user) 
        c.Set("user_id", user.ID)
        c.Set("email", claims["email"])
        c.Set("roles", claims["roles"])
        c.Next()
    }
}

func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRoles, exists := c.Get("roles")
        if !exists {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "status":  "error",
                "message": "Access denied - no roles information",
            })
            return
        }

        rolesInterface, ok := userRoles.([]interface{})
        if !ok {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "status":  "error",
                "message": "Invalid roles format",
            })
            return
        }

        var roles []string
        for _, role := range rolesInterface {
            if r, ok := role.(string); ok {
                roles = append(roles, r)
            }
        }

        for _, requiredRole := range requiredRoles {
            for _, userRole := range roles {
                if userRole == requiredRole {
                    c.Next()
                    return
                }
            }
        }

        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "status":  "error",
            "message": "Insufficient permissions",
            "details":"You are not authorized to access this endpoint.",
        })
    }
}

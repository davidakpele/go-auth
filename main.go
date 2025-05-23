package main

import (
	"api-service/config"
	"api-service/controllers"
	"api-service/db"
	"api-service/repositories"
	"api-service/services"
	"api-service/migrations"
	"api-service/routers"
	"log"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // Load configuration
    cfg := config.LoadConfig()

    // Connect to the database
    database, err := db.ConnectDatabase(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Migrate model
    if err := migrations.MigrateModels(database); err != nil {
        log.Fatalf("Database migration failed: %v", err)
    }

    // Create router
    router := gin.Default()

	// Custom CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
    router.Static("/static", "./static")  
    // Initialize dependencies

	authRepo := repositories.NewAuthRepository(database)
    userRepo := repositories.NewUserRepository(database)
    apiRepo := repositories.NewAPIRepository(database)


    apiService := services.NewAPIService(*apiRepo)
    authService := services.NewAuthService(*authRepo)
    userService := services.NewUserService(*userRepo)
    
    apiController := controllers.NewAPIController(*apiService)
    authController := controllers.NewAuthController(*authService)
    userController := controllers.NewUserController(*userService)
	
    // Register all routes by passing the router and dependencies
    routers.RegisterRoutes(router, authController, userController, apiController, userRepo) 

    // Start the server
    if err := router.Run(":7099"); err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
    gin.SetMode(gin.ReleaseMode)
}
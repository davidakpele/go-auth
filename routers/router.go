package routers

import (
	"api-service/controllers"
	"api-service/middleware"
	"api-service/repositories"
	"net/http"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes initializes the routes for the application
func RegisterRoutes(router *gin.Engine, 
	authController *controllers.AuthController, 
	userController *controllers.UserController,
	apiController *controllers.APIController,
	userRepo *repositories.UserRepository) {
	
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Access Denied",
			"status": "error", 
			"title": "Authentication Error", 
			"message": "Authorization Access",
			"details": "Something went wrong with authentication", 
			"code": "generic_authentication_error",
		})
        c.Abort()
	})

	// Handle undefined routes
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "EndPoint Not Found",
			"message": "Access Denied",
			"status": "error", 
			"title": "Authentication Error", 
			"details": "Something went wrong with authentication", 
			"code": "generic_authentication_error",
		})
	})

 	public := router.Group("/auth")
	{
		public.POST("/login", authController.Login)
		public.POST("/register", authController.Register)
		public.GET("/logout", authController.Logout)
		public.POST("/verify-account", authController.VerifyAccount)
		public.POST("/resend-otp", authController.ResendOTP)
	}

	// Private User routes
	private := router.Group("/api")
	private.Use(middleware.AuthenticationMiddleware(userRepo))
	private.Use(middleware.RoleMiddleware("USER", "ADMIN", "SUPER_USER",))
	{
		private.GET("/user/:id", userController.GetUserByID) 
		private.DELETE("/user/:id", userController.Delete) 
	}

	// Testing all users endpoints
	private_all_users_api_route := router.Group("/api/collection")
	private_all_users_api_route.Use(middleware.AuthenticationMiddleware(userRepo))
	private_all_users_api_route.Use(middleware.RoleMiddleware("ADMIN", "USER", "SUPER_USER", "MANAGER", "EDITOR"))
	{
		private_all_users_api_route.GET("", apiController.Collect)
		private_all_users_api_route.GET("/", apiController.Collect)
	}

}
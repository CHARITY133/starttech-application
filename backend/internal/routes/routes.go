package routes

import (
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Innocent9712/much-to-do/Server/MuchToDo/docs"
	"github.com/Innocent9712/much-to-do/Server/MuchToDo/internal/handlers"
	_ "github.com/Innocent9712/much-to-do/Server/MuchToDo/docs"
)

// RegisterRoutes sets up all application routes.
func RegisterRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	todoHandler *handlers.TodoHandler,
	healthHandler *handlers.HealthHandler,
	authMiddleware gin.HandlerFunc,
) {
	// Swagger documentation route (left outside /api, purely for docs UI)
	router.GET("/swagger/*any", func(c *gin.Context) {
		scheme := "http"
		if c.Request.TLS != nil || strings.HasPrefix(c.Request.Header.Get("X-Forwarded-Proto"), "https") {
			scheme = "https"
		}
		docs.SwaggerInfo.Host = c.Request.Host
		docs.SwaggerInfo.Schemes = []string{scheme}
		ginSwagger.WrapHandler(swaggerFiles.Handler)(c)
	})

	// Everything the frontend calls lives under /api, matching the
	// CloudFront /api/* -> ALB-Backend cache behavior.
	api := router.Group("/api")
	{
		// Public routes
		api.GET("/health", healthHandler.CheckHealth)

		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", userHandler.Register)
			authRoutes.POST("/login", userHandler.Login)
			authRoutes.POST("/logout", userHandler.Logout)
			authRoutes.GET("/username-check/:username", userHandler.CheckUsernameAvailability)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(authMiddleware)
		{
			taskRoutes := protected.Group("/tasks")
			{
				taskRoutes.POST("", todoHandler.CreateTodo)
				taskRoutes.GET("", todoHandler.GetAllTodos)
				taskRoutes.GET("/:id", todoHandler.GetTodoByID)
				taskRoutes.PUT("/:id", todoHandler.UpdateTodo)
				taskRoutes.DELETE("/:id", todoHandler.DeleteTodo)
			}

			userRoutes := protected.Group("/users")
			{
				userRoutes.GET("/me", userHandler.GetCurrentUser)
				userRoutes.PUT("/me", userHandler.UpdateUser)
				userRoutes.PUT("/me/password", userHandler.ChangePassword)
				userRoutes.DELETE("/me", userHandler.DeleteUser)
			}
		}
	}
}

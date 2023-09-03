package router

import (
	"task-5-pbi-btpns-daffajatmiko/controllers"
	"task-5-pbi-btpns-daffajatmiko/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// InitRoutes initializes the API routes and returns a Gin engine.
func InitRoutes(db *gorm.DB) *gin.Engine {
	// Create a new Gin router with default middleware
	router := gin.Default()

	// Middleware to set the database connection as a context variable
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	// User Routes
	router.POST("/users/login", controllers.Login)           // Route for user login
	router.POST("/users/register", controllers.CreateUser)    // Route for user registration
	router.PUT("/users/:userId", controllers.UpdateUser)      // Route to update user information
	router.DELETE("/users/:userId", controllers.DeleteUser)   // Route to delete a user account

	router.GET("/photos", controllers.GetPhoto)               // Route to retrieve photos

	// Middlewares for photo related routes
	authorized := router.Group("/").Use(middlewares.AuthMiddleware()) // Group of routes requiring authentication
	{
		authorized.POST("/photos", controllers.CreatePhoto)         // Route to create a new photo (authentication required)
		authorized.PUT("/photos/:photoId", controllers.UpdatePhoto)  // Route to update a photo (authentication required)
		authorized.DELETE("/photos/:photoId", controllers.DeletePhoto) // Route to delete a photo (authentication required)
	}

	return router
}

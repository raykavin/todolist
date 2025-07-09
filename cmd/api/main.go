package main

import (
	"ecommerce/internal/http/handler"
	repo "ecommerce/internal/infrastructure/repository/gorm"
	auth_uc "ecommerce/internal/service/auth"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize the database
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Initialize repositories
	userRepo := repo.NewUserRepository(db)

	// Initialize use cases
	authUsecase := auth_uc.NewAuthUsecase(userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
	}

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

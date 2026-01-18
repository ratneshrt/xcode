package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ratneshrt/xcode/database"
	"github.com/ratneshrt/xcode/handlers"
	"github.com/ratneshrt/xcode/models"
)

// postgresql://postgres:mysecretpassword@localhost:5432/postgres

func main() {
	database.Connect()
	if err := database.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	// publicRoutes := r.Group("/p")
	// {
	// 	publicRoutes.POST("/login", handlers.Login)
	// 	publicRoutes.POST("/register", handlers.Register)
	// }

	r.POST("/login", handlers.Login)
	r.POST("/register", handlers.Register)

	r.Run(":8080")
}

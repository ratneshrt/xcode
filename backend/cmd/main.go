package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ratneshrt/xcode/database"
	"github.com/ratneshrt/xcode/handlers"
	adminHandlers "github.com/ratneshrt/xcode/handlers/admin"
	"github.com/ratneshrt/xcode/middleware"
	"github.com/ratneshrt/xcode/models"
)

// postgresql://postgres:mysecretpassword@localhost:5432/postgres

func main() {
	database.ConnectAuthDB()
	database.ConnectProblemDB()

	if err := database.AuthDB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal(err)
	}

	if err := database.ProblemDB.AutoMigrate(
		&models.Problem{},
		&models.ProblemExample{},
		&models.ProblemTestCase{},
	); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.POST("/login", handlers.Login)
	r.POST("/register", handlers.Register)

	admin := r.Group("/admin")
	admin.Use(middleware.Authentication(), middleware.AdminOnly())
	{
		admin.POST("/problems", adminHandlers.CreateProblem)
	}

	r.Run(":8080")
}

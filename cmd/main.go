package main

import (
	"log"
	"net/http"

	"fumies/api/routes"
	"fumies/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.DatabaseConnection())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, "boop")
	})

	routes.UsePerfumeRoutes(r)
	routes.UseReviewRoutes(r)

	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	r := setupRouter()

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

package routes

import (
	"fumies/api/handlers"
	"fumies/middleware"

	"github.com/gin-gonic/gin"
)

func UseReviewRoutes(r *gin.Engine) {
	reviewGroup := r.Group("/review")
	reviewGroup.Use(middleware.JWTValidator())
	reviewGroup.GET("", handlers.GetReview)
	reviewGroup.POST("", handlers.CreateReview)
	reviewGroup.PUT("/:id", handlers.UpdateReview)
}

package routes

import (
	"fumies/api/handlers"
	"fumies/middleware"

	"github.com/gin-gonic/gin"
)

func UsePerfumeRoutes(r *gin.Engine) {
	perfumeGroup := r.Group("/perfume")
	perfumeGroup.Use(middleware.JWTValidator())
	perfumeGroup.GET("", handlers.GetPerfume)
	perfumeGroup.POST("", handlers.CreateOrUpdatePerfume)
	perfumeGroup.PUT("/:id", handlers.CreateOrUpdatePerfume)
	perfumeGroup.POST("/:id/wear", handlers.WearPerfume)
	perfumeGroup.POST("/recommend", handlers.RecommendPerfume)
}

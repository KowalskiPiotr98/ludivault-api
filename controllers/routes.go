package controllers

import "github.com/gin-gonic/gin"

func SetRoutes(r *gin.RouterGroup) {
	// platforms API
	platforms := r.Group("/platforms")
	platforms.GET("/", getPlatforms)
	platforms.GET("/:id", getPlatform)
	platforms.POST("/", createPlatform)
	platforms.PUT("/:id", updatePlatform)
	platforms.DELETE("/:id", deletePlatform)
}

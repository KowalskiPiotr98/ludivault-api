package controllers

import "github.com/gin-gonic/gin"

func SetRoutes(r *gin.RouterGroup) {
	//auth API
	auth := r.Group("/auth")
	auth.GET("", initAuth)
	auth.GET("/callback", authCallback)
	auth.POST("/logout", logout)

	// platforms API
	platforms := r.Group("/platforms")
	platforms.GET("", getPlatforms)
	platforms.GET("/:id", getPlatform)
	platforms.POST("", createPlatform)
	platforms.PUT("/:id", updatePlatform)
	platforms.DELETE("/:id", deletePlatform)

	// games API
	games := r.Group("/games")
	games.GET("", getGames)
	games.GET("/:id", getGame)
	games.GET("/:id/playthroughs", getPlaythroughsForGame)
	games.POST("", createGame)
	games.PUT("/:id", updateGame)
	games.DELETE("/:id", deleteGame)

	// playthroughs API
	playthroughs := r.Group("/playthroughs")
	playthroughs.GET("", getPlaythroughs)
	playthroughs.GET("/:id", getPlaythrough)
	playthroughs.POST("", createPlaythrough)
	playthroughs.PUT("/:id", updatePlaythrough)
	playthroughs.DELETE("/:id", deletePlaythrough)
}

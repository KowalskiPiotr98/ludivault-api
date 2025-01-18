package controllers

import (
	"github.com/KowalskiPiotr98/ludivault/auth"
	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.RouterGroup) {
	//auths API
	auths := r.Group("/auth")
	auths.GET("", initAuth)
	auths.GET("/callback", authCallback)
	auths.GET("/logout", logout)
	auths.GET("/providers", getProviders)
	auths.GET("/me", getUser)

	// platforms API
	platforms := r.Group("/platforms")
	platforms.Use(auth.GetLoginRequiredMiddleware())
	platforms.GET("", getPlatforms)
	platforms.GET("/:id", getPlatform)
	platforms.POST("", createPlatform)
	platforms.PUT("/:id", updatePlatform)
	platforms.DELETE("/:id", deletePlatform)

	// games API
	games := r.Group("/games")
	games.Use(auth.GetLoginRequiredMiddleware())
	games.GET("", getGames)
	games.GET("/:id", getGame)
	games.GET("/:id/playthroughs", getPlaythroughsForGame)
	games.GET("/:id/notes", getNoteTitles)
	games.POST("", createGame)
	games.PUT("/:id", updateGame)
	games.DELETE("/:id", deleteGame)

	// playthroughs API
	playthroughs := r.Group("/playthroughs")
	playthroughs.Use(auth.GetLoginRequiredMiddleware())
	playthroughs.GET("", getPlaythroughs)
	playthroughs.GET("/:id", getPlaythrough)
	playthroughs.POST("", createPlaythrough)
	playthroughs.PUT("/:id", updatePlaythrough)
	playthroughs.DELETE("/:id", deletePlaythrough)

	// notes API
	notes := r.Group("/notes")
	notes.Use(auth.GetLoginRequiredMiddleware())
	notes.GET("/:id", getNote)
	notes.POST("", createNote)
	notes.PUT("/:id", updateNote)
	notes.DELETE("/:id", deleteNote)
}

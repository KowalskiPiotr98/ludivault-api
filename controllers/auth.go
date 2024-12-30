package controllers

import (
	"github.com/KowalskiPiotr98/ludivault/auth"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func initAuth(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		gothic.BeginAuthHandler(c.Writer, c.Request)
		return
	}

	if err = initUserSession(&user, c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func authCallback(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Warnf("Error while authenticating user: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = initUserSession(&user, c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func logout(c *gin.Context) {
	if err := auth.RemoveUserSession(c); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

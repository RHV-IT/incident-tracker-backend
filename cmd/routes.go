package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *application) routes() http.Handler {
	g := gin.Default()
	g.RedirectTrailingSlash = true

	v1 := g.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", a.authMiddleware(), a.register)
			authGroup.POST("/login", a.login)
			authGroup.PUT("/update", a.authMiddleware(), a.update)
			authGroup.PUT("/disable", a.authMiddleware(), a.disable)
			authGroup.PUT("/enable",  a.authMiddleware(), a.enable)
		}

	}

	return  g
}
package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (a *application) routes() http.Handler {
	g := gin.Default()
	g.RedirectTrailingSlash = true

	allowedOrigins := strings.Split(a.origins, ",")

	for i, origin := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(origin)
	}

	g.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
			authGroup.PUT("/enable", a.authMiddleware(), a.enable)
			authGroup.PUT("/resetpassword", a.authMiddleware(), a.resetPassword)
			authGroup.PUT("/userResetPassword", a.authMiddleware(), a.userResetPassword)
		}
		v1.POST("/incidents", a.reportIncident)
		v1.GET("/incidents", a.authMiddleware(), a.getIncidents)
		v1.GET("/users", a.authMiddleware(), a.getUsers)
		v1.GET("/user", a.authMiddleware(), a.getUser)
		v1.PATCH("/incidents/:id/status", a.authMiddleware(), a.updateIncidentStatus)
		v1.POST("/incidents/:id/management", a.authMiddleware(), a.submitIncidentManagement)
		v1.GET("/incidents/:id/management", a.authMiddleware(), a.getIncidentManagement)
		v1.PUT("/incidents/:id/management", a.authMiddleware(), a.updateIncidentManagement)
		v1.GET("/incidents/:id/managementlogs", a.authMiddleware(), a.getIncidentLogs)
		v1.POST("/incidents/comments", a.authMiddleware(), a.addComment)
		v1.GET("/incidents/comments", a.authMiddleware(), a.getComments)
		v1.GET("/searchUsers", a.authMiddleware(), a.searchUsers)
		v1.POST("/deathreport", a.deathReport)
		v1.PUT("/deathreport/:id", a.updateDeathReport)
		v1.GET("/deathreport/:query", a.searchDeathReport)
	}

	return g
}

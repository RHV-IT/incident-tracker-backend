package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *application) routes() http.Handler {
	g := gin.Default()

	return  g
}
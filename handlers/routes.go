package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", HomeHandler())
	r.GET("/screen", ScreenHandler())
	r.GET("/ws", WebSocketHandler())
}
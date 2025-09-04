package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/", HomeHandler(db))
	r.GET("/screen", ScreenHandler(db))
	r.GET("/ws", WebSocketHandler(db))
}
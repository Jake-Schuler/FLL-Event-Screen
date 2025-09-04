package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ScreenHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "screen.tmpl", gin.H{
			"title": "FLL Event Screen - Public Screen",
		})
	}
}

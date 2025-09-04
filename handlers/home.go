package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HomeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "index.tmpl", gin.H{
			"title": "FLL Event Screen - Dashboard",
		})
	}
}

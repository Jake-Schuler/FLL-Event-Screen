package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func ScreenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(200, "screen.tmpl", gin.H{
			"title":      "FLL Event Screen - Public Screen",
			"event_name": os.Getenv("EVENT_NAME"),
		})
	}
}

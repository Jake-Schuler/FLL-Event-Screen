package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func ScreenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		view := c.Query("view")
		switch view {
		case "queue":
			c.HTML(200, "upcoming_queue.tmpl", gin.H{
				"title":      "FLL Event Screen - Queue Screen",
				"event_name": os.Getenv("EVENT_NAME"),
			})
		default:
			c.HTML(200, "screen.tmpl", gin.H{
				"title":      "FLL Event Screen - Public Screen",
				"event_name": os.Getenv("EVENT_NAME"),
			})
		}
	}
}

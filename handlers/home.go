package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func HomeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		view := c.Query("view")
		switch view {
		case "timer":
			c.HTML(200, "timer_control.tmpl", gin.H{
				"title": "FLL Event Screen - Timer Controller",
				"event_name": os.Getenv("EVENT_NAME"),
			})
		default:
			c.HTML(200, "index.tmpl", gin.H{
				"title":      "FLL Event Screen - Dashboard",
				"event_name": os.Getenv("EVENT_NAME"),
			})
		}
	}
}

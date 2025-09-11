package handlers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func ScreenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		view := c.Query("view")
		backgroundColor := c.Query("background")
		if backgroundColor == "" {
			switch view {
			case "table":
				backgroundColor = "0a0f3c"
			default:
				backgroundColor = "FF00FF"
			}
		}
		switch view {
		case "table":
			c.HTML(200, "table_screen.tmpl", gin.H{
				"title":            "FLL Event Screen - Table Screen",
				"event_name":       os.Getenv("EVENT_NAME"),
				"background_color": backgroundColor,
			})
		default:
			c.HTML(200, "screen.tmpl", gin.H{
				"title":            "FLL Event Screen - Public Screen",
				"event_name":       os.Getenv("EVENT_NAME"),
				"background_color": backgroundColor,
			})
		}
	}
}

package codingame

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)


func Ping() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
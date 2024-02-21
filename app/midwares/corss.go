package midwares

import (
	"github.com/gin-gonic/gin"
)

func Corss(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Authorization")
	c.Writer.Header().Set("Access-Control-Max-Age", "172800")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
		return
	}
	c.Next()
}
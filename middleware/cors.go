package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 解决跨域问题
func Cors(c *gin.Context) {
	method := c.Request.Method
	if method != "" {

		c.Header("Access-Control-Allow-Origin", c.GetHeader("origin"))

		c.Header("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT")

		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")

		c.Header("Access-Control-Allow-Credentials", "true")

	}

	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
}

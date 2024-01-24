package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, UPDATE, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Accept-Encoding, Authorization, Connection, Content-Length, Content-Type, Host, User-Agent, Token, Access-Control-Request-Method, Access-Control-Request-Headers")
		// c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "9000")
		c.Set("content-type", "application/json")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

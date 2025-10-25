package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging(c *gin.Context) {
	start := time.Now()

	c.Next()

	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}

	key := slog.Attr{}

	switch {
	case strings.HasSuffix(path, ":id"):
		key.Key = "user_id"
		key.Value = slog.StringValue(c.Param("id"))
	case strings.HasSuffix(path, ":username"):
		key.Key = "username"
		key.Value = slog.StringValue(c.Param("username"))
	}

	slog.Info("Incoming request:",
		"http.server.request.duration", time.Since(start).String(),
		"http.request.method", c.Request.Method,
		"http.response.status_code", c.Writer.Status(),
		"http.route", path,
		"server.address", c.Request.URL.Path,
		"http.request.host", c.Request.Host,
		key,
	)
}

func APIKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		got := c.GetHeader("X-API-Key")

		if got == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing api key",
			})
			c.Abort()
			return
		}

		if got != key {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

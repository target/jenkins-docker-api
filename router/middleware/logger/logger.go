// Package logger provides log handling using logrus package.
//
// Based on github.com/gin-gonic/contrib/ginrus but adds more options.
package logger

import (
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// LogHandler returns a gin.HandlerFunc (middleware) that logs requests using logrus.
//
// Requests with errors are logged using logrus.Error().
// Requests without errors are logged using logrus.Info().
//
// It receives:
//   1. A time package format string (e.g. time.RFC3339).
//   2. A boolean stating whether to use UTC time zone or local.
func LogHandler(logger *logrus.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		// prevent us from logging the k8s health checks every 10s
		if strings.Contains(path, "health") == false {
			entry := logger.WithFields(logrus.Fields{
				"api-id":      c.GetHeader("J-Api-Id"),
				"api-version": c.GetHeader("J-Api-Version"),
				"body":        c.Value("payload"),
				"status":      c.Writer.Status(),
				"method":      c.Request.Method,
				"path":        path,
				"ip":          c.ClientIP(),
				"latency":     latency,
				"user-agent":  c.Request.UserAgent(),
				"time":        end.Format(timeFormat),
			})

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				entry.Error(strings.TrimSpace(c.Errors.String()))
			} else {
				entry.Info()
			}
		}

	}
}

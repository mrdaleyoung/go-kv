package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

// responseWriter wraps gin.ResponseWriter to capture the status code
type responseWriter struct {
	gin.ResponseWriter
	status int
}

// WriteHeader captures the status code
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs the details of each request
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		//Log some data points
		log.Printf("Request: %s %s | Status: %d | Duration: %v | Client IP: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			c.ClientIP(),
		)
	}
}

// SecureMiddleware applies security settings
func SecureMiddleware() gin.HandlerFunc {
	options := secure.Options{
		STSSeconds:            31536000,
		ContentTypeNosniff:    true,
		ContentSecurityPolicy: "default-src 'self; frame-ancestors 'none';'",
		ReferrerPolicy:        "no-referrer",        // Referrer policy
		FeaturePolicy:         "geolocation 'none'", // Feature policy
	}

	secureMiddleware := secure.New(options)

	return func(c *gin.Context) {
		writer := &responseWriter{ResponseWriter: c.Writer}

		// Call the secure middleware
		secureMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Next()
		})).ServeHTTP(writer, c.Request)

		// Set the response status
		c.Status(writer.status)
	}
}

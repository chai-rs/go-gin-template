package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
)

// RateLimitMiddleware returns a middleware that limits the number of requests per user.
func RateLimitMiddleware(limiter *limiter.Limiter) gin.HandlerFunc {
	return mgin.NewMiddleware(limiter)
}

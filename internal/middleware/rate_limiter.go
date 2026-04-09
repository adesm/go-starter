package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)
		
		ctx := c.Request.Context()
		
		// Fixed Window Algorithm
		// Increment the counter for this IP
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			// If Redis is down, we might want to log it and let the request pass 
			// or fail closed. Here we fail open but log.
			c.Next()
			return
		}

		// If this is the first request in the window, set expiration
		if count == 1 {
			rdb.Expire(ctx, key, time.Minute)
		}

		// Limit to 100 requests per minute
		if count > 100 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false, 
				"error": "Rate limit exceeded. Please try again in a minute.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

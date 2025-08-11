package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type bucket struct {
	mu   sync.Mutex
	hits []time.Time
}

var (
	rlStore = struct {
		mu sync.Mutex
		m  map[string]*bucket
	}{m: make(map[string]*bucket)}
)

// RateLimit limits requests per key (ip+path) to max hits within window duration
func RateLimit(max int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + c.FullPath()
		rlStore.mu.Lock()
		b, ok := rlStore.m[key]
		if !ok {
			b = &bucket{}
			rlStore.m[key] = b
		}
		rlStore.mu.Unlock()

		now := time.Now()
		cutoff := now.Add(-window)

		b.mu.Lock()
		// purge old
		i := 0
		for _, t := range b.hits {
			if t.After(cutoff) {
				b.hits[i] = t
				i++
			}
		}
		b.hits = b.hits[:i]
		if len(b.hits) >= max {
			b.mu.Unlock()
			c.Header("Retry-After", window.String())
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		b.hits = append(b.hits, now)
		b.mu.Unlock()

		c.Next()
	}
}

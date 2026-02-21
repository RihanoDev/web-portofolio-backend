package middleware

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type encodeResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w encodeResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func EncodeResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Only encode responses for API routes, skip websockets and uploads
		if !strings.HasPrefix(path, "/api/v1") ||
			strings.HasSuffix(path, "/upload") ||
			strings.Contains(path, "/ws") {
			c.Next()
			return
		}

		w := &encodeResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = w

		c.Next()

		if w.Status() == 204 || w.body.Len() == 0 {
			return
		}

		encodedBytes := []byte(base64.StdEncoding.EncodeToString(w.body.Bytes()))

		w.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(encodedBytes)))
		w.ResponseWriter.Header().Set("Content-Type", "text/plain")
		w.ResponseWriter.Header().Set("X-Encoded-Response", "true")

		w.ResponseWriter.Write(encodedBytes)
	}
}

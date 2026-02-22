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
	body       *bytes.Buffer
	statusCode int
	done       bool // set true setelah encode selesai ditulis
}

func (w *encodeResponseWriter) Write(b []byte) (int, error) {
	if w.done {
		return w.ResponseWriter.Write(b)
	}
	return w.body.Write(b)
}

func (w *encodeResponseWriter) WriteHeader(statusCode int) {
	if w.done {
		w.ResponseWriter.WriteHeader(statusCode)
		return
	}
	// Simpan statusCode untuk nanti — jangan tulis sekarang
	w.statusCode = statusCode
}

func (w *encodeResponseWriter) WriteHeaderNow() {
	// Cegah flush prematur selama body masih dikumpulkan
	// (Gin memanggil ini dari c.Status, tapi kita belum siap)
	if w.done {
		w.ResponseWriter.WriteHeaderNow()
	}
}

func (w *encodeResponseWriter) Status() int {
	if w.statusCode != 0 {
		return w.statusCode
	}
	return w.ResponseWriter.Status()
}

func (w *encodeResponseWriter) Written() bool {
	return w.done
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
			statusCode:     200,
		}
		c.Writer = w

		c.Next()

		statusCode := w.statusCode
		if statusCode == 0 {
			statusCode = 200
		}

		if statusCode == 204 || w.body.Len() == 0 {
			w.done = true
			w.ResponseWriter.WriteHeader(statusCode)
			return
		}

		encodedBytes := []byte(base64.StdEncoding.EncodeToString(w.body.Bytes()))

		header := w.ResponseWriter.Header()
		header.Set("Content-Length", strconv.Itoa(len(encodedBytes)))
		header.Set("Content-Type", "text/plain; charset=utf-8")
		header.Set("X-Encoded-Response", "true")

		w.done = true
		w.ResponseWriter.WriteHeader(statusCode)
		w.ResponseWriter.Write(encodedBytes) //nolint:errcheck
	}
}

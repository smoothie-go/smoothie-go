package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// stub
func Render(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming unsupported"})
		return
	}

	for progress := 0; progress <= 100; progress += 5 {
		_, err := fmt.Fprintf(c.Writer, "data: {\"progress\": %d}\n\n", progress)
		if err != nil {
			return
		}
		flusher.Flush()
		time.Sleep(500 * time.Millisecond)
	}
}

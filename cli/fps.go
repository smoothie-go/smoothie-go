package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"
)

func GetFramerate(filePath string) (float64, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	metadata, err := ffprobe.ProbeURL(ctx, filePath)
	if err != nil {
		return 0, fmt.Errorf("error getting metadata: %w", err)
	}

	for _, stream := range metadata.Streams {
		if stream.CodecType == "video" {
			return parseFramerate(stream.RFrameRate)
		}
	}

	return 0, fmt.Errorf("no video stream found in the file")
}

func parseFramerate(rFrameRate string) (float64, error) {
	parts := strings.Split(rFrameRate, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid framerate format: %s", rFrameRate)
	}
	numerator, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid numerator in framerate: %w", err)
	}
	denominator, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid denominator in framerate: %w", err)
	}
	if denominator == 0 {
		return 0, fmt.Errorf("denominator in framerate cannot be zero")
	}
	return float64(numerator) / float64(denominator), nil
}

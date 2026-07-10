package services

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"time"
)

func VideoProccessor(dir string) int {
	//Navigating to folder where video is there
	fmt.Printf("compressing ...")

	os.MkdirAll(dir+"/streaming_output", os.ModePerm)

	args := []string{
		"-i", "original.mp4", // Input file
		"-threads", "0", // Max CPU cores
		"-preset", "ultrafast", // Max encoding speed
		"-vcodec", "libx264", // H.264 video codec
		"-vf", "scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2",
		"-crf", "23", // Compression value
		"-g", "60",
		"-keyint_min", "60",
		"-sc_threshold", "0", // Making a key frame every 60 frames
		"-f", "hls", // For output to be in hls format
		"-hls_time", "2", // Making 2 seconds segments
		"-hls_playlist_type", "event",
		"-hls_segment_filename", "streaming_output/file%03d.ts",
		"streaming_output/stream.m3u8",
	}

	cmd := exec.Command("ffmpeg", args...)
	cmd.Dir = dir

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Failed to create stderr pipe: %v\n", err)
		return 0
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Start()
	if err != nil {
		fmt.Printf("FFmpeg failed: %v\nLogs: %s\n", err, stderr.String())
	}

	timeRegex := regexp.MustCompile(`time=(\d{2}:\d{2}:\d{2}\.\d{2})`)
	var lastTimestamp string

	// Split by '\r' because FFmpeg uses carriage returns to refresh the same line
	reader := bufio.NewReader(stderrPipe)
	for {
		line, readErr := reader.ReadString('\r')

		// Simultaneously write everything to your original log buffer
		stderr.WriteString(line)

		if readErr != nil {
			if readErr != io.EOF {
				fmt.Printf("Error reading stream: %v\n", readErr)
			}
			break // End of process stream
		}

		// Look for the "time=00:00:00.00" progress updates
		matches := timeRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			lastTimestamp = matches[1]
		}
	}

	// 5. Wait for FFmpeg to completely finish execution
	err = cmd.Wait()
	if err != nil {
		// Your exact original error logging format
		fmt.Printf("FFmpeg failed: %v\nLogs: %s\n", err, stderr.String())
		return 0
	}

	// 6. Convert the very last captured timestamp to total seconds
	totalSeconds := parseTimestampToSeconds(lastTimestamp)
	return totalSeconds
}

// Helper function to turn "00:02:15.50" into an integer (135 seconds)
func parseTimestampToSeconds(timestamp string) int {
	if timestamp == "" {
		return 0
	}
	t, err := time.Parse("15:04:05.99", timestamp)
	if err != nil {
		return 0
	}

	// Math conversion: (hours * 3600) + (minutes * 60) + seconds
	return (t.Hour() * 3600) + (t.Minute() * 60) + t.Second()
}

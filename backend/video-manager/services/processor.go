package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func VideoProccessor(dir string) {
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

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.Dir = dir
	err := cmd.Run()

	if err != nil {
		fmt.Printf("FFmpeg failed: %v\nLogs: %s\n", err, stderr.String())
	}
}

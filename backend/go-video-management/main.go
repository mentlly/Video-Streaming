package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	port := os.Getenv("GO_VIDEO_MANAGEMENT_PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "Ok"}`)
	})

	//Creating upload folder to store uploaded videos
	os.MkdirAll("./uploads", os.ModePerm)

	http.HandleFunc("/api/upload", uploadVideoHandler)

	log.Printf("Go Video Management Server booting up internally on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start Go server: %v", err)
	}
}

const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 * 1024

func uploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "File to large or bad request", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}

	fmt.Printf("Uploaded File: %s\n", handler.Filename)
	fmt.Printf("File Size: %d bytes\n", handler.Size)

	//Creating directory for uploaded video to write in
	videoId := generateVideoId()
	dir := "./uploads/" + videoId
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	dst, err := os.Create(filepath.Join(dir, "original.mp4"))
	if err != nil {
		http.Error(w, "Internal server error creating file", http.StatusBadRequest)
		return
	}
	defer dst.Close()

	//Copying the file to the created filename
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusBadRequest)
		return
	}

	videoProccessor(dir)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Video uploaded successfully"))
}

// currently not used to generate a unique id
func videoHasher(file multipart.File) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	hashInBytes := hasher.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(hashInBytes), nil
}

// Genrates a random string of length 10 for videoId
func generateVideoId() string {
	alphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	length := 10
	max := big.NewInt(62)
	rstr := ""

	for i := 1; i <= length; i++ {
		secureNum, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}

		rstr += string(alphabet[int(secureNum.Int64())])
	}
	return rstr
}

func videoProccessor(dir string) {
	//Navigating to folder where video is there
	fmt.Printf("compressing ...")

	os.MkdirAll(dir+"/streaming_output", os.ModePerm)

	args := []string{
		"-i", "original.mp4", // Input file
		"-vcodec", "libx264", // H.264 video codec
		"-vf", "scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2",
		"-crf", "23", // Compression value
		"-g", "60",
		"-keyint_min", "60",
		"-sc_threshold", "0", // Making a key frame every 60 frames
		"-f", "hls", // For output to be in hls format
		"-hls_time", "5", // Making 5 seconds segments
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

package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"video-manager/services"
	"video-manager/utils"
)

const MAX_UPLOAD_SIZE = 2 * 1024 * 1024 * 1024

func UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("context: %v\n", r.Context().Value("userId"))

	account_id, ok := r.Context().Value("userId").(int)
	if !ok || account_id <= 0 {
		// If the cast fails or it's empty, the middleware didn't set it properly
		http.Error(w, "Unauthorized: Missing identity profile", http.StatusUnauthorized)
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

	title := r.FormValue("title")
	if strings.TrimSpace(title) == "" {
		title = handler.Filename
	}
	description := r.FormValue("description")

	//Creating directory for uploaded video to write in
	videoId := utils.GenerateVideoId()
	dir := "./uploads/" + videoId
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		http.Error(w, "Internal server error creating object folder", http.StatusBadRequest)
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
	}

	duration := services.VideoProccessor(dir)

	services.UploadVideoDb(account_id, videoId, title, description, duration)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Video uploaded successfully"))
}

func GetLatestVideoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	videos, err := services.GetLatestVideos(1)
	if err != nil {
		http.Error(w, "Error fetching videos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(videos)
}

func GetIdVideoHandler(w http.ResponseWriter, r *http.Request) {
	videoId := r.PathValue("VideoId")

	// 1. Target the exact folder where this specific video's chunks live
	dirPath := fmt.Sprintf("./uploads/%s/streaming_output", videoId)

	// 2. Double-check if the folder actually exists to avoid confusing errors
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// 3. Apply necessary streaming headers before handing off to Go's core file engine
	filename := r.PathValue("filename")

	// 💡 FIX 1: Prevent Directory Traversal & Folder Browsing
	// If 'filename' is empty, or if it points to a directory path instead of a file, block it!
	cleanedFile := filepath.Clean(filename)
	if filename == "" || strings.HasSuffix(r.URL.Path, "/") || cleanedFile == "." || cleanedFile == ".." {
		http.Error(w, "Access Forbidden", http.StatusForbidden)
		return
	}

	fullDiskPath := filepath.Join(dirPath, cleanedFile)

	// 💡 FIX 2: Explicitly check that the requested item is a FILE, not a folder on disk
	fileInfo, err := os.Stat(fullDiskPath)
	if os.IsNotExist(err) || fileInfo.IsDir() {
		http.NotFound(w, r)
		return
	}

	if strings.HasSuffix(filename, ".m3u8") {
		w.Header().Set("Content-Type", "application/x-mpegURL")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	} else if strings.HasSuffix(filename, ".ts") {
		w.Header().Set("Content-Type", "video/MP2T")
	}

	// 4. Strip the URL structure up to this folder context and serve
	prefix := fmt.Sprintf("/api/video/get/%s", videoId)
	handler := http.StripPrefix(prefix, http.FileServer(http.Dir(dirPath)))

	handler.ServeHTTP(w, r)
}

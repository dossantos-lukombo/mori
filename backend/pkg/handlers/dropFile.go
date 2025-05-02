package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"mori/pkg/utils"
)

// Define the upload path – adjust as needed.
const uploadPath = "./fileUploads"

// isPathSafe checks if the given path is safe and within the upload directory
func isPathSafe(path string) bool {
	// Get absolute paths
	absUploadPath, err := filepath.Abs(uploadPath)
	if err != nil {
		return false
	}
	absFilePath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Check if the file path is within the upload directory
	return strings.HasPrefix(absFilePath, absUploadPath)
}

// sanitizeFilename removes potentially dangerous characters from the filename
func sanitizeFilename(filename string) string {
	// Remove any path separators
	filename = filepath.Base(filename)
	
	// Remove any null bytes
	filename = strings.ReplaceAll(filename, "\x00", "")
	
	// Remove any potentially dangerous characters
	dangerous := []string{"..", "~", "/", "\\"}
	for _, d := range dangerous {
		filename = strings.ReplaceAll(filename, d, "")
	}
	
	return filename
}

// UploadFiles handles file uploads.
// It expects a multipart form with one or more files under the key "files".
func (h *Handler) UploadFiles(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)
	// Allow only POST requests.
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ensure the upload directory exists.
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
			http.Error(w, "Error creating upload directory", http.StatusInternalServerError)
			return
		}
	}

	// Parse the multipart form (limit: 20MB).
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Retrieve files with key "files".
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Sanitize the filename
		sanitizedFilename := sanitizeFilename(fileHeader.Filename)
		if sanitizedFilename == "" {
			http.Error(w, "Invalid filename", http.StatusBadRequest)
			return
		}

		// Create destination file path
		dstPath := filepath.Join(uploadPath, sanitizedFilename)

		// Verify the path is safe
		if !isPathSafe(dstPath) {
			http.Error(w, "Invalid file path", http.StatusBadRequest)
			return
		}

		// Create the file
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file data to the destination.
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "File(s) uploaded successfully")
}

// ListFiles returns a JSON array of objects, each containing the file name, its size, and upload date.
func (h *Handler) ListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(uploadPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Folder does not exist: treat as empty.
			files = []os.DirEntry{}
		} else {
			http.Error(w, "Error reading upload directory", http.StatusInternalServerError)
			return
		}
	}

	// Define a struct to hold file information.
	type FileInfo struct {
		Name       string `json:"name"`
		Size       int64  `json:"size"`
		UploadDate string `json:"uploadDate"`
	}

	var filesInfo []FileInfo
	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				// Skip files with errors getting info.
				continue
			}
			filesInfo = append(filesInfo, FileInfo{
				Name:       file.Name(),
				Size:       info.Size(),
				UploadDate: info.ModTime().Format("2006-01-02 15:04:05"), // Format as desired.
			})
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filesInfo)
}

// DeleteFile deletes a file from the upload folder.
// It expects the URL pattern: /api/files/{filename}
func (h *Handler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	w = utils.ConfigHeader(w)

	// Handle preflight OPTIONS request.
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Filename not specified", http.StatusBadRequest)
		return
	}
	filename := parts[3]
	filePath := filepath.Join(uploadPath, filename)

	// If the file doesn't exist, return success.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "File not found (already deleted)")
		return
	}

	if err := os.Remove(filePath); err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "File deleted successfully")
}

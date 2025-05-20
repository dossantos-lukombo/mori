package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Path to default image location
const defaultImage = "frontend/public/images/default.svg"

// Upload directory for user avatars
const avatarUploadDir = "imageUploads"

// Initialize function will run when the package is imported
func init() {
	// Ensure imageUploads directory exists
	if err := os.MkdirAll(avatarUploadDir, 0755); err != nil {
		fmt.Printf("Warning: Failed to create avatar upload directory: %v\n", err)
	}

	// Ensure the default image directory exists
	defaultImageDir := filepath.Dir(defaultImage)
	if err := os.MkdirAll(defaultImageDir, 0755); err != nil {
		fmt.Printf("Warning: Failed to create default image directory: %v\n", err)
	}

	// Check if default image exists, if not create a simple SVG
	if _, err := os.Stat(defaultImage); os.IsNotExist(err) {
		fmt.Println("Default avatar not found, creating a simple one...")
		defaultSVG := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" width="100" height="100">
			<defs>
				<linearGradient id="bgGradient" x1="0%" y1="0%" x2="100%" y2="100%">
					<stop offset="0%" stop-color="#8B5CF6" />
					<stop offset="100%" stop-color="#1F2937" />
				</linearGradient>
			</defs>
			<rect width="100" height="100" rx="50" ry="50" fill="url(#bgGradient)" />
			<circle cx="50" cy="35" r="25" fill="#FFFFFF" />
			<circle cx="50" cy="100" r="40" fill="#FFFFFF" />
		</svg>`
		
		if err := os.WriteFile(defaultImage, []byte(defaultSVG), 0644); err != nil {
			fmt.Printf("Warning: Failed to create default avatar: %v\n", err)
		} else {
			fmt.Println("Default avatar created successfully.")
		}
	}
}

// Creates new file and reads image bytes into it
// returns path to new image
// returns default avatar or provided one
func SaveAvatar(r *http.Request) string {
	// Read data from request
	file, fileHeader, errRead := r.FormFile("avatar")
	if errRead != nil {
		return defaultImage
	}
	defer file.Close()
	
	// get content type -> png, gif or jpeg
	contentType := fileHeader.Header["Content-Type"][0]
	
	// Ensure the upload directory exists
	if err := os.MkdirAll(avatarUploadDir, 0755); err != nil {
		fmt.Printf("Failed to create upload directory: %v\n", err)
		return defaultImage
	}
	
	// Create new file with custom name instead of temp file
	extension := getExtensionFromContentType(contentType)
	if extension == "" {
		return defaultImage
	}
	
	// Create a unique filename with timestamp
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%s/avatar_%d%s", avatarUploadDir, timestamp, extension)
	
	// Create the file
	localFile, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return defaultImage
	}
	defer localFile.Close()

	// read data in new file
	fileData, err := io.ReadAll(file)
	if err != nil {
		return defaultImage
	}
	localFile.Write(fileData)
	return strings.Replace(localFile.Name(), "\\", "/", -1)
}

/* ------------------------- for posts and comments ------------------------- */
// Creates new file and reads image bytes into it
// returns path to new image
// returns no path if not exist
func SaveImage(r *http.Request) string {
	// Read data from request
	file, fileHeader, errRead := r.FormFile("image")
	if errRead != nil {
		return ""
	}
	defer file.Close()
	
	// get content type -> png, gif or jpeg
	contentType := fileHeader.Header["Content-Type"][0]
	
	// Ensure the upload directory exists
	if err := os.MkdirAll(avatarUploadDir, 0755); err != nil {
		fmt.Printf("Failed to create upload directory: %v\n", err)
		return ""
	}
	
	// Create new file with custom name instead of temp file
	extension := getExtensionFromContentType(contentType)
	if extension == "" {
		return ""
	}
	
	// Create a unique filename with timestamp
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%s/image_%d%s", avatarUploadDir, timestamp, extension)
	
	// Create the file
	localFile, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		return ""
	}
	defer localFile.Close()

	// read data in new file
	fileData, err := io.ReadAll(file)
	if err != nil {
		return ""
	}
	localFile.Write(fileData)
	return strings.Replace(localFile.Name(), "\\", "/", -1)
}

// Helper function to get file extension from content type
func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	default:
		return ""
	}
}

// This function is kept for backward compatibility if needed
func createTempFile(fileType string) (*os.File, error) {
	var localFile *os.File
	var err error

	// Make sure the directory exists
	if err := os.MkdirAll(avatarUploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	extension := getExtensionFromContentType(fileType)
	if extension == "" {
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}

	// Use CreateTemp with the new directory
	localFile, err = os.CreateTemp(avatarUploadDir, fmt.Sprintf("*%s", extension))
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	
	return localFile, nil
}

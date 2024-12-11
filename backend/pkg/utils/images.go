package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Patht to default image location
const defaultImage = "imageUpload/default.svg"

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
	// Create empty local file with correct file extension
	localFile, err := createTempFile(contentType)
	// If type not recognized return default image path
	if err != nil {
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
	// Create empty local file with correct file extension
	localFile, err := createTempFile(contentType)
	// If type not ecognized return default image path
	if err != nil {
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

// creates empty local file based on file type
func createTempFile(fileType string) (*os.File, error) {
	tempDir, err := os.MkdirTemp("", "imageUpload")
	if err != nil {
		return nil, err
	}

	var pattern string
	switch fileType {
	case "image/jpeg":
		pattern = "*.jpg"
	case "image/png":
		pattern = "*.png"
	case "image/gif":
		pattern = "*.gif"
	default:
		return nil, fmt.Errorf("type de fichier non supporté : %s", fileType)
	}

	// Crée un fichier temporaire dans le répertoire spécifié
	localFile, err := os.CreateTemp(tempDir, pattern)
	if err != nil {
		return nil, err
	}

	return localFile, nil
}

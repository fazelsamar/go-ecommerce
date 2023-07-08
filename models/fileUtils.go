package models

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var uploadPath string = "mediafiles"

// SaveUploadedFile saves the uploaded image file and returns the file path or an error
func saveUploadedFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Validate file MIME type
	if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		return "", errors.New("invalid file type. Only images are allowed")
	}

	// Create a unique file name
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))

	// Specify the path to save the file
	savePath := filepath.Join(uploadPath, filename)

	// Create the directory if it doesn't exist
	if err := createDirectory(uploadPath); err != nil {
		return "", err
	}

	// Create the destination file
	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	// Convert the file path separators to forward slashes
	savePath = filepath.ToSlash(savePath)

	return savePath, nil
}

// Helper function to create the directory if it doesn't exist
func createDirectory(dirPath string) error {
	// Check if the directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Create the directory with 0755 permissions
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to generate a unique file path for the uploaded file
// func GenerateFilePath(filename string) string {
// 	// Get the current working directory
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		return ""
// 	}

// 	// Create a new directory uploadPath if it doesn't exist
// 	uploadsDir := filepath.Join(cwd, uploadPath)
// 	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
// 		os.Mkdir(uploadsDir, 0755)
// 	}

// 	// Generate a unique file name using a timestamp and the original file name
// 	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), filename)

// 	// Construct the full file path by combining the uploads directory and the unique file name
// 	filePath := filepath.Join(uploadsDir, fileName)

// 	return "/" + filepath.ToSlash(filePath)
// }

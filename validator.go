package uploader

import (
	"errors"
	"io"
	"net/http"
)

var ErrInvalidFileType = errors.New("invalid picture file")

var ErrFileTooLarge = errors.New("file size exceeds the limit")

func ValidateFileSize(file io.Seeker, maxSize int64) error {
	// Get the current position of the file pointer
	currPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	// Seek to the end to get the total size of the file
	fileSize, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	// Reset the file pointer back to the original position
	if _, err := file.Seek(currPos, io.SeekStart); err != nil {
		return err
	}

	// Check if the file size exceeds the allowed limit
	if fileSize > maxSize {
		return ErrFileTooLarge
	}

	return nil
}

func ValidateFileType(file io.ReadSeeker, fileTypes ...string) (string, error) {
	// Validate the MIME type (only allow jpg and png)
	fileHeader := make([]byte, 512) // Read the first 512 bytes for MIME type detection
	_, err := file.Read(fileHeader)
	if err != nil {
		return "", err
	}

	// Reset the file pointer to the beginning after reading
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	// Detect file MIME type
	mimeType := http.DetectContentType(fileHeader)
	if !InArrayStr(fileTypes, mimeType) {
		return "", ErrInvalidFileType
	}

	return mimeType, nil
}

// Contains checks if a slice contains a specific string.
func InArrayStr(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

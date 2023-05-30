package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsImageFile(t *testing.T) {
	// Test with valid image file extensions
	validExtensions := []string{".jpg", ".jpeg", ".JPG", ".JPEG"}
	for _, ext := range validExtensions {
		assert.True(t, IsImageFile("example"+ext), "Expected true, got false for extension %s", ext)
	}

	// Test with invalid file extensions
	invalidExtensions := []string{".png", ".gif", ".bmp", ".txt"}
	for _, ext := range invalidExtensions {
		assert.False(t, IsImageFile("example"+ext), "Expected false, got true for extension %s", ext)
	}

	// Test with file names that have no extensions
	noExtensionFileNames := []string{"example", "image123", "picture-2021"}
	for _, fileName := range noExtensionFileNames {
		assert.False(t, IsImageFile(fileName), "Expected false, got true for file name %s", fileName)
	}
}

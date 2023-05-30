package utils

import (
	"path/filepath"
	"strings"
)

func IsImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg"
}

package s3

import (
	"io"
	"os"
)

type S3Client interface {
	ListObjects() ([]Object, error)
	GetObject(key string) (io.ReadCloser, error)
	PutObject(key string, file *os.File) error
}

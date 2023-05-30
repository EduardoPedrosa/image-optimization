package main

import (
	"log"
	"os"
	"strconv"

	"github.com/eduardopedrosa/image-optimization/processor"
	"github.com/eduardopedrosa/image-optimization/s3"
	"github.com/eduardopedrosa/image-optimization/utils"
	"github.com/joho/godotenv"
)

var (
	maxImageSize int
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables:", err)
	}

	maxImageSize, err = strconv.Atoi(os.Getenv("IMAGE_MAX_SIZE"))

	if err != nil {
		log.Println("Error during conversion of IMAGE_MAX_SIZE using 102400 instead")
		maxImageSize = 102400
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_KEY")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("S3_BUCKET")

	s3Client, err := s3.NewS3Client(accessKey, secretKey, region, bucket)

	if err != nil {
		log.Fatal("Error creating connection with S3 Client")
	}

	// Channel to receive errors from goroutines
	errChan := make(chan error)

	// Channel to synchronize goroutine completion
	doneChan := make(chan struct{})

	// Image counter
	var numImages int

	// List objects in the image folder
	objects, err := s3Client.ListObjects()
	if err != nil {
		log.Fatal("Could not get Objects from S3")
	}

	for _, obj := range objects {
		if utils.IsImageFile(obj.Name) && obj.Size > int64(maxImageSize) {
			numImages++
			go func(key string) {
				p := processor.NewProcessor(
					s3Client,
					maxImageSize,
					processor.ProcessorOptions{})

				err = p.Process(key)
				if err != nil {
					errChan <- err
					numImages--
					return
				}
				log.Println("Optimized image:", key)
				doneChan <- struct{}{}
			}(obj.Name)
		}
	}

	// Wait for goroutines to complete
	go func() {
		for i := 0; i < numImages; i++ {
			<-doneChan
		}
		close(errChan)
	}()

	// Check for errors in goroutines
	for err := range errChan {
		log.Println("Error during image processing:", err)
	}

	log.Println("Image processing completed")
}

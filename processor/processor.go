package processor

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/disintegration/imaging"
	"github.com/eduardopedrosa/image-optimization/s3"
)

type Processor struct {
	s3                       s3.S3Client
	maxImageSize             int
	localTempDir             string
	localTempPath            string
	defaultImageWidth        int
	defaultImageHeight       int
	initialQualityPercentage int
}

type ProcessorOptions struct {
	LocalTempDir       string
	DefaultImageWidth  int
	DefaultImageHeight int
}

const (
	initialQualityPercentage = 80
)

func NewProcessor(s3Client s3.S3Client, maxImageSize int, options ProcessorOptions) *Processor {
	localTempDir := "/tmp/image-optimization"
	if options.LocalTempDir != "" {
		localTempDir = options.LocalTempDir
	}
	defaultImageWidth := 1280
	if options.DefaultImageWidth != 0 {
		defaultImageWidth = options.DefaultImageWidth
	}
	defaultImageHeight := 720
	if options.DefaultImageHeight != 0 {
		defaultImageHeight = options.DefaultImageHeight
	}

	return &Processor{
		s3:                       s3Client,
		maxImageSize:             maxImageSize,
		localTempDir:             localTempDir,
		defaultImageWidth:        defaultImageWidth,
		defaultImageHeight:       defaultImageHeight,
		initialQualityPercentage: initialQualityPercentage,
	}
}

func (p *Processor) Process(key string) error {
	// Download the image from S3
	downloadErr := p.downloadImage(key)
	if downloadErr != nil {
		return downloadErr
	}

	defer p.cleanUpTempData(key)

	// Optimize the image
	optimizeErr := p.optimizeImage(key)
	if optimizeErr != nil {
		return optimizeErr
	}

	// Replace the original image in S3
	replaceErr := p.replaceImage(key)
	if replaceErr != nil {
		return replaceErr
	}

	// Signal that the goroutine has completed
	return nil
}

func (p *Processor) downloadImage(key string) error {
	p.ensureTempDirectoryExists()

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	p.localTempPath = p.localTempDir + "/" + timestamp + ".jpg"

	file, err := os.Create(p.localTempPath)
	if err != nil {
		return err
	}
	defer file.Close()

	obj, err := p.s3.GetObject(key)

	if err != nil {
		return err
	}

	_, err = io.Copy(file, obj)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) optimizeImage(key string) error {
	// Open the original image
	img, err := imaging.Open(p.localTempPath)
	if err != nil {
		return err
	}

	// Resize the image defaultImageWidth x defaultImageHeight
	resizedImg := imaging.Fit(img, p.defaultImageWidth, p.defaultImageHeight, imaging.Lanczos)

	// Compress the image to a maximum of maxImageSize
	buf := &bytes.Buffer{}
	err = imaging.Encode(buf, resizedImg, imaging.JPEG, imaging.JPEGQuality(initialQualityPercentage))

	if err != nil {
		return err
	}

	encodedImg := buf.Bytes()

	// Check the size of the compressed image
	if len(encodedImg) > p.maxImageSize {
		// If the image still exceeds maxImageSize, try reducing the quality further
		for quality := initialQualityPercentage; len(encodedImg) > p.maxImageSize && quality > 0; quality -= 5 {
			buf := &bytes.Buffer{}
			err = imaging.Encode(buf, resizedImg, imaging.JPEG, imaging.JPEGQuality(quality))
			if err != nil {
				return err
			}

			encodedImg = buf.Bytes()
		}
	}

	// Save the optimized image
	outputFile, err := os.Create(p.localTempPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write the image data to the output file
	_, err = outputFile.Write(encodedImg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) replaceImage(key string) error {
	file, err := os.Open(p.localTempPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = p.s3.PutObject(key, file)

	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) ensureTempDirectoryExists() error {

	_, err := os.Stat(p.localTempDir)
	if err == nil {
		// The parent directory already exists
		return nil
	}

	// Check if the error is because the directory does not exist
	if os.IsNotExist(err) {
		// Create the parent directory with read, write, and execute permissions for the owner
		err = os.MkdirAll(p.localTempDir, 0700)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) cleanUpTempData(key string) error {
	err := os.RemoveAll(p.localTempPath)
	if err != nil {
		return err
	}

	return nil
}

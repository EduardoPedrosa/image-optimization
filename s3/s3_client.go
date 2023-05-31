package s3

import (
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type AwsS3Client struct {
	svc    s3iface.S3API
	bucket string
}

type Object struct {
	Name string
	Size int64
}

func NewS3Client(accessKey, secretKey, region, bucket string) (*AwsS3Client, error) {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		log.Fatal("Error creating AWS SDK session:", err)
		return nil, err
	}

	return &AwsS3Client{
		svc:    s3.New(sess),
		bucket: bucket,
	}, nil
}

func (s3Client *AwsS3Client) ListObjects() ([]Object, error) {
	var objects []Object

	imageFolderPath := os.Getenv("S3_PREFIX")

	err := s3Client.svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(s3Client.bucket),
		Prefix: aws.String(imageFolderPath),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			objects = append(objects, Object{Name: *obj.Key, Size: *obj.Size})
		}
		return true
	})

	if err != nil {
		return nil, err
	}

	return objects, nil
}

func (s3Client *AwsS3Client) GetObject(key string) (io.ReadCloser, error) {
	objOutput, err := s3Client.svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s3Client.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return objOutput.Body, nil
}

func (s3Client *AwsS3Client) PutObject(key string, file *os.File) error {
	_, err := s3Client.svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Client.bucket),
		Key:    aws.String(key),
		Body:   file,
	})

	if err != nil {
		return err
	}

	return nil
}

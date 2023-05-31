package s3

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockS3Client struct {
	mock.Mock
	s3iface.S3API
}

// ListObjectsV2Pages is a mock implementation of the ListObjectsV2Pages method of the S3API interface.
func (m *MockS3Client) ListObjectsV2Pages(input *s3.ListObjectsV2Input, fn func(*s3.ListObjectsV2Output, bool) bool) error {
	args := m.Called(input, fn)
	return args.Error(0)
}

// GetObject is a mock implementation of the GetObject method of the S3API interface.
func (m *MockS3Client) GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

// PutObject is a mock implementation of the PutObject method of the S3API interface.
func (m *MockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func TestNewS3Client(t *testing.T) {
	accessKey := "your-access-key"
	secretKey := "your-secret-key"
	region := "your-region"
	bucket := "your-bucket"

	s3Client, err := NewS3Client(accessKey, secretKey, region, bucket)

	// Verify that no error occurred while creating the client
	assert.Nil(t, err)

	// Verify that the client was created correctly
	assert.NotNil(t, s3Client)
}

func TestNewS3Client_Error(t *testing.T) {
	accessKey := "your-access-key"
	secretKey := "your-secret-key"
	region := "your-region"
	bucket := "your-bucket"

	s3Client, err := NewS3Client(accessKey, secretKey, region, bucket)

	// Verify that no error occurred while creating the client
	assert.Nil(t, err)

	// Verify that the client was created correctly
	assert.NotNil(t, s3Client)
}

func TestListObjects(t *testing.T) {
	// Create a MockS3Client object
	mockS3Client := new(MockS3Client)

	// Set the expected behavior for the ListObjectsV2Pages method
	mockS3Client.On("ListObjectsV2Pages", mock.Anything, mock.Anything).
		Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(*s3.ListObjectsV2Output, bool) bool)
		// Simulate the callback function call
		fn(&s3.ListObjectsV2Output{
			Contents: []*s3.Object{
				{
					Key:  aws.String("object-key-1"),
					Size: aws.Int64(100),
				},
				{
					Key:  aws.String("object-key-2"),
					Size: aws.Int64(200),
				},
			},
		}, true)
	})

	// Create an AwsS3Client object with the mocked S3 client
	s3Client := &AwsS3Client{
		svc:    mockS3Client,
		bucket: "your-bucket",
	}

	// Execute the method to be tested
	objects, err := s3Client.ListObjects()

	// Check for no errors
	assert.Nil(t, err)

	// Check if the returned objects are as expected
	expectedObjects := []Object{
		{Name: "object-key-1", Size: 100},
		{Name: "object-key-2", Size: 200},
	}
	assert.Equal(t, expectedObjects, objects)

	// Check if the ListObjectsV2Pages method was called correctly
	mockS3Client.AssertCalled(t, "ListObjectsV2Pages", &s3.ListObjectsV2Input{
		Bucket: aws.String("your-bucket"),
		Prefix: aws.String(""),
	}, mock.Anything)
}

func TestListObjects_Error(t *testing.T) {
	// Create a MockS3Client object
	mockS3Client := new(MockS3Client)

	// Set the expected behavior for the ListObjectsV2Pages method to return an error
	expectedError := errors.New("list objects error")
	mockS3Client.On("ListObjectsV2Pages", mock.Anything, mock.Anything).
		Return(expectedError)

	// Create an AwsS3Client object with the mocked S3 client
	s3Client := &AwsS3Client{
		svc:    mockS3Client,
		bucket: "your-bucket",
	}

	// Execute the method to be tested
	objects, err := s3Client.ListObjects()

	// Check if the error is as expected
	assert.Equal(t, expectedError, err)

	// Check that the returned objects slice is nil
	assert.Nil(t, objects)
}

func TestGetObject(t *testing.T) {
	// Create a mock S3 client
	mockS3Client := new(MockS3Client)

	// Define the expected behavior for the GetObject method
	expectedOutput := &s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewBufferString("test content")),
	}
	mockS3Client.On("GetObject", mock.Anything).Return(expectedOutput, nil)

	// Create an AwsS3Client object with the mocked S3 client
	s3Client := &AwsS3Client{
		svc:    mockS3Client,
		bucket: "your-bucket",
	}

	// Execute the method to be tested
	objectContent, err := s3Client.GetObject("test-object-key")

	// Verify that there are no errors
	assert.Nil(t, err)

	// Read the object content
	content, _ := io.ReadAll(objectContent)

	// Verify that the content is as expected
	assert.Equal(t, []byte("test content"), content)

	// Verify that the GetObject method was called with the expected input
	mockS3Client.AssertCalled(t, "GetObject", &s3.GetObjectInput{
		Bucket: aws.String("your-bucket"),
		Key:    aws.String("test-object-key"),
	})
}

func TestGetObject_Error(t *testing.T) {
	// Create a mock S3 client
	mockS3Client := new(MockS3Client)

	// Define the expected behavior for the GetObject method to return an error
	expectedError := errors.New("get object error")
	mockS3Client.On("GetObject", mock.Anything).Return(nil, expectedError)

	// Create an AwsS3Client object with the mocked S3 client
	s3Client := &AwsS3Client{
		svc:    mockS3Client,
		bucket: "your-bucket",
	}

	// Execute the method to be tested
	objectContent, err := s3Client.GetObject("test-object-key")

	// Verify that the error is as expected
	assert.Equal(t, expectedError, err)

	// Verify that the objectContent is nil
	assert.Nil(t, objectContent)
}

func TestPutObject(t *testing.T) {
	// Create a mock S3 client
	mockS3Client := new(MockS3Client)

	// Define the expected behavior for the PutObject method
	mockS3Client.On("PutObject", mock.Anything).Return(&s3.PutObjectOutput{}, nil)

	// Create an AwsS3Client object with the mocked S3 client
	s3Client := &AwsS3Client{
		svc:    mockS3Client,
		bucket: "your-bucket",
	}

	// Open a test file
	file, _ := os.Open("test-file.txt")
	defer file.Close()

	// Execute the method to be tested
	err := s3Client.PutObject("test-object-key", file)

	// Verify that there are no errors
	assert.Nil(t, err)

	// Verify that the PutObject method was called with the expected input
	mockS3Client.AssertCalled(t, "PutObject", &s3.PutObjectInput{
		Bucket: aws.String("your-bucket"),
		Key:    aws.String("test-object-key"),
		Body:   file,
	})
}

func TestPutObject_Error(t *testing.T) {
	// Create a mock S3 client
	mockS3Client := new(MockS3Client)

	// Define the expected behavior for the PutObject method to return an error
	expectedError := errors.New("put object error")
	mockS3Client.On("PutObject", mock.Anything).Return(nil, expectedError)

	// Create an AwsS3Client object with the mocked S3 client
	s3Client := &AwsS3Client{
		svc:    mockS3Client,
		bucket: "your-bucket",
	}

	// Open a test file
	file, _ := os.Open("test-file.txt")
	defer file.Close()

	// Execute the method to be tested
	err := s3Client.PutObject("test-object-key", file)

	// Verify that the error is as expected
	assert.Equal(t, expectedError, err)
}

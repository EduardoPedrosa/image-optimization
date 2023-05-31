# Image Optimization

This project is a multithreaded script written in Go that retrieves JPG images from Amazon S3 and optimizes them by reducing their size to a desired dimension. It aims to automate the process of optimizing images on a large scale.

## Project Structure

The project consists of the following files and directories:

- `main.go`: The entry point of the program that initiates the image processing.
- `s3/`: A package that encapsulates the communication with Amazon S3.
- `processor/`: A package that handles image processing.
- `utils/`: A package containing utility functions.

## Technologies Used

- Language: Go (Golang)
- Libraries:

  - `github.com/aws/aws-sdk-go v1.44.271`: The official AWS SDK for Go, used for interacting with Amazon S3.
  - `github.com/disintegration/imaging v1.6.2`: A powerful imaging library for Go, used for image resizing and optimization.
  - `github.com/jmespath/go-jmespath v0.4.0`: A Go implementation of JMESPath, used for querying JSON data.
  - `github.com/joho/godotenv v1.5.1`: A Go library for loading environment variables from a `.env` file.

## Running the Project

To run the project, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/EduardoPedrosa/image-optimization.git
   ```

2. Navigate to the project directory:

   ```bash
   cd image-optimization
   ```

3. Create a `.env` file and set the required environment variables (e.g., AWS access key, secret key, region, S3 bucket, etc.).

4. Build and run the project using the following command:

   ```bash
   go run main.go
   ```

   Note: Make sure you have Go installed and configured on your machine.
   

## Testing the Project

To test the project, run this command:

   ```bash
      go test ./... -cover
   ```
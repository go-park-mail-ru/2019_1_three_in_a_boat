package handlers

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"image"
	"image/jpeg"
	"io"
)

const (
	imageFolder = "img/"
)

// The basket that we use
var (
	bucketName = "hexagon-game"
	region     = "eu-north-1"
)

// Singletons
// Initialized once in init() function
var uploader *s3manager.Uploader
var svc *s3.S3

// Initializing an uploader with the session and default options
func init() {
	// The session the S3 Uploader will use
	// and that the SDK will use to load
	// credentials from environment variables
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	// Create an uploader with the session and default options
	uploader = s3manager.NewUploader(sess)

	// Initialize instance of the S3 client with a session
	// Needed for init BatchDeleteIterator
	svc = s3.New(sess)
}

// Convert golang image.Image to bytes
// Needed for initializing UploadInput struct
func imageToBuffer(img image.Image) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func saveObject(keyName string, r io.Reader) error {
	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    aws.String(imageFolder + keyName),
		Body:   r,
	}

	// Perform an upload.
	_, err := uploader.Upload(upParams)
	if err != nil {
		return err
	}

	return nil
}

// Save image in AWS bucket
func SaveImage(img image.Image, name string) error {
	buf, err := imageToBuffer(img)
	if err != nil {
		return err
	}

	err = saveObject(name, buf)
	if err != nil {
		return err
	}

	return nil
}

// Delete image from AWS bucket
func DeleteImage(img image.Image, name string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(imageFolder + name),
	}

	_, err := svc.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr
		} else {
			return err
		}
	}

	return nil
}

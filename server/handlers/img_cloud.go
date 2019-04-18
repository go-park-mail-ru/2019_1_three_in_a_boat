package handlers

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"image"
	"io"
)

var (
	bucketName = "hexagon-images"
	region = "eu-central-1"
)

func imageToBuffer(img image.Image) *bytes.Buffer {
	var b []byte
	return bytes.NewBuffer(b)
}

func saveObject(keyName string, r io.Reader) error {
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
	}))
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)
	// Upload input parameters
	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &keyName,
		Body:   r,
	}
	// Perform an upload.
	_, err := uploader.Upload(upParams)
	if err != nil {
		return err
	}

	return nil
}

func SaveImage(img image.Image, name string) error {
	buf := imageToBuffer(img)
	err := saveObject(name, buf)
	if err != nil {
		return err
	}
	return nil
}
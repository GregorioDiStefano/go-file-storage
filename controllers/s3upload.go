package controller

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (upload Upload) doActualUpload(data io.ReadCloser, key string, fn string) error {
	//TODO:why no recovery when crash here? ex: failed DNS lookup.

	s3uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(upload.awsRegion)}))

	_, err := s3uploader.Upload(&s3manager.UploadInput{
		Body:   data,
		Bucket: aws.String(upload.awsBucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", key, fn)),
	})

	if err != nil {
		return err
	}

	return nil
}

func (upload Upload) doActualDelete(deleteKey, fileKey, filename string) error {
	return nil
}

package controller

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/GregorioDiStefano/go-file-storage/utils"
)

type S3Upload struct{}

func (s3up S3Upload) upload(data io.ReadCloser, key string, fn string) error {
	awsRegion := utils.Config.GetString("aws.region")
	awsBucket := utils.Config.GetString("aws.bucket")

	uploader := s3manager.NewUploader(session.New(&aws.Config{Region: aws.String(awsRegion)}))
	_, err := uploader.Upload(&s3manager.UploadInput{
		Body:   data,
		Bucket: aws.String(awsBucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", key, fn)),
	})

	if err != nil {
		return err
	}

	return nil
}

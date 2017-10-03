package amazon

import (
	"fmt"
	"io"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	log *logrus.Entry
}

func New(log *logrus.Entry) *S3 {
	return &S3{
		log: log.WithField("provider", "s3"),
	}
}

func GetManifest(manifestURL url.URL) (io.ReadCloser, error) {
	bucket := manifestURL.Host
	key := manifestURL.Path

	cfg := aws.NewConfig()
	awsSession := session.New(cfg)
	s3Service := s3.New(awsSession)

	result, err := s3Service.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("error getting s3 object '%s' in bucket '%s': %s", key, bucket, err)
	}
	return result.Body, nil
}

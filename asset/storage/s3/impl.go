package s3

import (
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/janoszen/openshiftci-inspector/asset/storage"
)

type s3AssetStorage struct {
	storage.AbstractAssetStorage

	s3     *s3.S3
	bucket string
}

func (s *s3AssetStorage) Store(jobID string, name string, mime string, data []byte) error {
	key := "/" + jobID + "/" + name
	_, err := s.s3.PutObject(
		&s3.PutObjectInput{
			ACL:           aws.String(s3.BucketCannedACLPublicRead),
			Body:          bytes.NewReader(data),
			Bucket:        aws.String(s.bucket),
			ContentLength: aws.Int64(int64(len(data))),
			ContentType:   aws.String(mime),
			Key:           aws.String(key),
		},
	)
	return err
}

func (s *s3AssetStorage) Fetch(jobID string, name string) (data []byte, err error) {
	key := "/" + jobID + "/" + name
	get, err := s.s3.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return
	}
	defer func() {
		_ = get.Body.Close()
	}()
	data, err = ioutil.ReadAll(get.Body)
	if err != nil {
		return
	}
	return
}
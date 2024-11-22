package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Event struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func FetchS3Object(event S3Event) (*s3.GetObjectOutput, error) {
	bucket, key := extractS3Details(event)

	sess := session.Must(session.NewSession())
	s3Client := s3.New(sess)
	return s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
}

func extractS3Details(event S3Event) (string, string) {
	return event.Records[0].S3.Bucket.Name, event.Records[0].S3.Object.Key
}

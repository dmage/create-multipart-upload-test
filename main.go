package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	bucket = "obulatov-create-multipart-upload-test"
)

func test() (*s3.CreateMultipartUploadOutput, error) {
	svc := s3.New(
		session.New(
			&aws.Config{
				LogLevel:         aws.LogLevel(aws.LogDebugWithRequestErrors),
				S3ForcePathStyle: aws.Bool(false),
			},
		),
	)

	input := &s3.CreateMultipartUploadInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String("file.txt"),
		ServerSideEncryption: aws.String(s3.ServerSideEncryptionAes256),
		ContentType:          aws.String("application/octet-stream"),
	}

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			log.Fatal(err)
		}
	}()

	return svc.CreateMultipartUpload(input)
}

func main() {
	result, err := test()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
}

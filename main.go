package main

import (
	"flag"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	bucket = "obulatov-create-multipart-upload-test"
	key    = "file.txt"
)

func test() error {
	svc := s3.New(
		session.New(
			&aws.Config{
				LogLevel:         aws.LogLevel(aws.LogDebugWithRequestErrors),
				S3ForcePathStyle: aws.Bool(false),
			},
		),
	)

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return err
	}
	defer func() {
		_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Body:        strings.NewReader("Hello, world"),
		Bucket:      aws.String(bucket),
		ContentType: aws.String("application/octet-stream"),
		Key:         aws.String(key),
	})
	if err != nil {
		return err
	}

	_, err = svc.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &s3.Delete{
			Objects: []*s3.ObjectIdentifier{
				{
					Key: aws.String(key),
				},
			},
		},
	})
	if err != nil {
		_, deleteErr := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
		log.Println(deleteErr)
	}
	return err
}

func main() {
	flag.StringVar(&bucket, "bucket", bucket, "a name of the test bucket")
	flag.StringVar(&key, "key", key, "a key for the test object")
	flag.Parse()

	err := test()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("SUCCESS!")
}

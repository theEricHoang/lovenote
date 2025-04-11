package imageservice

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Presigner struct {
	PresignClient *s3.PresignClient
}

func NewPresigner(presignClient *s3.PresignClient) *Presigner {
	return &Presigner{PresignClient: presignClient}
}

func (p *Presigner) PresignPut(ctx context.Context, key, contentType string) (string, error) {
	req, err := p.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String("lovenote-images"),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(15*time.Minute))

	if err != nil {
		log.Printf("Error presigning PUT url for %v: %v", key, err)
		return "", err
	}

	return req.URL, nil
}

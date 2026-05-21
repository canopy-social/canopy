package s3

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	sysconfig "github.com/sumi-devs/canopy-social/canopy/pkg/config"
)

type Client struct {
	client  *s3.Client
	bucket  string
	cdnBase string
}

func NewClient(cfg *sysconfig.Config) (*Client, error) {
	var opts []func(*config.LoadOptions) error

	opts = append(opts, config.WithRegion(cfg.S3.Region))

	if cfg.S3.AccessKey != "" && cfg.S3.SecretKey != "" {
		creds := credentials.NewStaticCredentialsProvider(cfg.S3.AccessKey, cfg.S3.SecretKey, "")
		opts = append(opts, config.WithCredentialsProvider(creds))
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load S3 config: %w", err)
	}

	s3Opts := func(o *s3.Options) {
		if cfg.S3.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.S3.Endpoint)
			o.UsePathStyle = true
		}
	}

	s3Client := s3.NewFromConfig(awsCfg, s3Opts)

	return &Client{
		client:  s3Client,
		bucket:  cfg.S3.Bucket,
		cdnBase: cfg.S3.CDNBase,
	}, nil
}

func (c *Client) Upload(ctx context.Context, key string, contentType string, reader io.Reader) (string, error) {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
		ACL:         s3types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to put S3 object: %w", err)
	}

	if c.cdnBase != "" {
		return fmt.Sprintf("%s/%s", c.cdnBase, key), nil
	}

	if c.client.Options().BaseEndpoint != nil && *c.client.Options().BaseEndpoint != "" {
		return fmt.Sprintf("%s/%s/%s", *c.client.Options().BaseEndpoint, c.bucket, key), nil
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", c.bucket, key), nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete S3 object: %w", err)
	}
	return nil
}

func (c *Client) Bucket() string {
	return c.bucket
}

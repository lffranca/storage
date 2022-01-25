package storage

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

func NewS3(ctx context.Context, bucket string) (Storage, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := new(clientS3)
	client.bucket = bucket
	client.client = s3.NewFromConfig(cfg)

	return client, nil
}

type clientS3 struct {
	bucket string
	client *s3.Client
}

func (pkg *clientS3) Download(ctx context.Context, key *string) (io.Reader, error) {
	buff := new(bytes.Buffer)
	buffW := NewBufferWriterAt(buff)

	downloader := manager.NewDownloader(pkg.client)
	if _, err := downloader.Download(ctx, buffW, &s3.GetObjectInput{
		Bucket: &pkg.bucket,
		Key:    key,
	}); err != nil {
		return nil, err
	}

	return buff, nil
}

func (pkg *clientS3) PreSign(ctx context.Context, key, contentType *string) (*string, error) {
	psClient := s3.NewPresignClient(pkg.client)
	res, err := psClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:              &pkg.bucket,
		Key:                 key,
		ResponseContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	return &res.URL, nil
}

func (pkg *clientS3) ListObjects(ctx context.Context) ([]*File, error) {
	output, err := pkg.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &pkg.bucket,
	})
	if err != nil {
		return nil, err
	}

	var files []*File
	for _, object := range output.Contents {
		size := int(object.Size)
		files = append(files, &File{
			Key:          object.Key,
			Size:         &size,
			LastModified: object.LastModified,
		})
	}

	return files, nil
}

func (pkg *clientS3) Upload(ctx context.Context, key, contentType *string, data io.Reader) (err error) {
	uploader := manager.NewUploader(pkg.client)
	if _, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Key:         key,
		ContentType: contentType,
		Body:        data,
		Bucket:      &pkg.bucket,
	}); err != nil {
		return err
	}

	return
}

package awsS3

import (
	"axis/ecommerce-backend/configs"
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/disintegration/imaging"
	"image/png"
	"io"
	"mime/multipart"
)

type Client struct {
	ses      *session.Session
	uploader *s3manager.Uploader
}

func (c Client) CreateImage(file *multipart.FileHeader, width int, height int, style string, id uint) (io.Reader, string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, "", err
	}
	defer src.Close()

	buffer := make([]byte, file.Size)
	_, err = src.Read(buffer)
	if err != nil {
		return nil, "", err
	}
	bb := bytes.NewReader(buffer)

	fileName := fmt.Sprintf("/part_images/%d/file.%s.png", id, style)
	newImage, err := png.Decode(bb)
	if err != nil {
		return nil, "", err
	}
	buff := new(bytes.Buffer)
	resized := imaging.Resize(newImage, width, height, imaging.Lanczos)
	err = png.Encode(buff, resized)
	if err != nil {
		return nil, "", err
	}

	return bytes.NewReader(buff.Bytes()), fileName, nil
}

func (c Client) UploadImage(fileName string, file io.Reader) (string, error) {
	res, err := c.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String("axisforestry"),
		Key:         aws.String(fileName),
		Body:        file,
		ACL:         aws.String("public-read"),
		ContentType: aws.String("image/png"),
	})

	if err != nil {
		return "", err
	}

	return res.Location, nil
}

func (c Client) DeleteResource(bucket, key string) error {
	svc := s3.New(c.ses)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return err
}

func NewS3Client() (*Client, error) {
	getAwsConfig, err := configs.NewAwsConfig()
	if err != nil {
		return nil, err
	}

	ses, err := session.NewSession(&aws.Config{
		Region:      aws.String(getAwsConfig.AwsRegion),
		Credentials: credentials.NewStaticCredentials(getAwsConfig.AwsKeyId, getAwsConfig.AwsSecretKeyId, ""),
	})

	if err != nil {
		return nil, err
	}

	uploader := s3manager.NewUploader(ses)
	return &Client{uploader: uploader, ses: ses}, nil
}

// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package uploader

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"git.tech.kora.id/go/utility/random"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Session session instances
var S3Session *session.Session

// S3Setup returns a new Session created from SDK defaults, config files,
// environment, and user provided config files. Once the Session is created
// it can be mutated to modify the Config or Handlers. The Session is safe to
// be read concurrently, but it should not be written to concurrently.
func S3Setup() (e error) {
	awsID := os.Getenv("AWS_KEY")
	awsSecret := os.Getenv("AWS_SECRET")
	awsRegion := os.Getenv("AWS_REGION")

	creds := credentials.NewStaticCredentials(awsID, awsSecret, "")
	if _, e = creds.Get(); e != nil {
		return
	}

	S3Session, e = session.NewSession(&aws.Config{
		Region:           aws.String(awsRegion),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
	})

	return
}

// UploadImage uploading image as filebuffer
func UploadImage(dir string, name string, fb []byte, thumb bool) (url string, e error) {
	// if name is empty we generate some random string
	if name == "" {
		name = random.String(15, random.Numeric)
	}

	var thumbImage *Image

	image := NewImage(fmt.Sprintf("%s/%s", dir, name), fb)

	if thumb {
		thumbImage, _ = image.MakeThumbnail()
	}

	// doing seamless upload so its not blocking the process
	go func() {
		if thumb {
			uploader(thumbImage.FileName(), thumbImage.Data)
		}

		uploader(image.FileName(), image.Data)
	}()

	url = getURL(image.FileName())

	if thumb {
		url = getURL(thumbImage.FileName())
	}

	return
}

func uploader(fn string, fb []byte) (url string, e error) {
	bucket := os.Getenv("AWS_BUCKET")
	ct := http.DetectContentType(fb)

	if _, e = s3.New(S3Session).PutObject(&s3.PutObjectInput{
		Bucket:          aws.String(bucket),
		Key:             aws.String(fn),
		Body:            bytes.NewReader(fb),
		ContentType:     aws.String(ct),
		ContentEncoding: aws.String("Base64"),
	}); e == nil {
		url = getURL(fn)
	}

	return
}

func getURL(res string) string {
	awsBucket := os.Getenv("AWS_BUCKET")
	awsRegion := os.Getenv("AWS_REGION")

	return fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", awsRegion, awsBucket, res)
}

func getExtention(ct string) (ext string) {
	switch ct {
	case "image/jpeg":
		ext = ".jpg"
		break
	case "image/png":
		ext = ".png"
		break
	case "image/gif":
		ext = ".gif"
		break
	}

	return
}

func getFilename(url string) string {
	return strings.Replace(url, getURL(""), "", 1)
}

// DownloadImage download image to base64
func DownloadImage(url string) (res string, imageBytes []byte) {
	bucket := os.Getenv("AWS_BUCKET")
	// replace dulu urlnya menjadi
	filename := getFilename(url)

	buff := &aws.WriteAtBuffer{}
	s3dl := s3manager.NewDownloader(S3Session)
	if _, err := s3dl.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}); err == nil {
		imageBytes = buff.Bytes()
		res = base64.StdEncoding.EncodeToString(imageBytes)
	}

	return
}

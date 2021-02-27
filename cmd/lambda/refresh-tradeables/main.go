package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/jlaffaye/ftp"

	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func putS3(r io.Reader, name string, bucket string, svc *s3.S3) (err error) {

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	rs := bytes.NewReader(buf)

	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Body:   rs,
		Key:    aws.String(name),
	}

	_, err = svc.PutObject(input)

	return err
}

func Handler() error {
	c, err := ftp.Dial("ftp.nasdaqtrader.com:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	err = c.Login("anonymous", "anonymous@domain.com")
	if err != nil {
		return err
	}

	toDownload := []string{
		"nasdaqlisted.txt",
		"otherlisted.txt",
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	S3Service := s3.New(sess)

	bucketName := os.Getenv("BUCKET")

	for _, f := range toDownload {
		r, err := c.Retr(fmt.Sprintf("SymbolDirectory/%s", f))
		if err != nil {
			return err
		}

		err = putS3(r, f, bucketName, S3Service)
		r.Close()
		if err != nil {
			return err
		}
	}

	err = c.Quit()
	return err

}

func main() {
	lambda.Start(Handler)
}

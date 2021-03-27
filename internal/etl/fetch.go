package etl

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	db "github.com/fischersean/phish-food/internal/database"
	"github.com/fischersean/phish-food/internal/reddit"
	"github.com/fischersean/phish-food/internal/stocks"

	"time"
)

func GetPosts(sub string, limit int, auth reddit.AuthToken, conn db.Connection) (p []reddit.Post, err error) {

	p, err = reddit.GetHot(sub, limit, auth)
	if err != nil {
		return p, err
	}

	for i := range p {
		p[i].DownloadTime = time.Now()
		p[i].Comments, err = reddit.FetchPostComments(p[i], auth)
		if err != nil {
			return p, err
		}
	}

	go ArchivePost(p, sub, conn)

	return p, err
}

func GetTradeableSecurities(svc *s3.S3, bucket string) ([]stocks.Stock, error) {

	input0 := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("otherlisted.txt"),
	}
	input1 := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String("nasdaqlisted.txt"),
	}

	result0, err := svc.GetObject(input0)
	if err != nil {
		return nil, err
	}
	defer result0.Body.Close()

	result1, err := svc.GetObject(input1)
	if err != nil {
		return nil, err
	}
	defer result1.Body.Close()

	s0, err := stocks.FromNasdaqOtherListed(result0.Body)
	if err != nil {
		return nil, err
	}

	s1, err := stocks.FromNasdaqListed(result1.Body)
	if err != nil {
		return nil, err
	}

	stocksPopulation := append(s0, s1...)

	return stocksPopulation, nil

}

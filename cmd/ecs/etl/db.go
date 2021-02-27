package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/fischersean/phish-food/internal/report"
	"github.com/fischersean/phish-food/internal/stocks"

	"github.com/fischersean/phish-food/internal/database"

	"time"
)

func getTradeableSecurities(svc *s3.S3, bucket string) ([]stocks.Stock, error) {

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

func putRecord(conn database.Connection, sr []report.StockReport, sub string, t time.Time) (err error) {

	e := database.NewEtlResultsRecord(sr, sub, t)
	err = conn.PutEtlResultsRecord(e)
	return err
}

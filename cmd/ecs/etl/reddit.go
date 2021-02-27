package main

import (
	"github.com/fischersean/phish-food/internal/database"
	"github.com/fischersean/phish-food/internal/reddit"
	"github.com/fischersean/phish-food/internal/report"
	"github.com/fischersean/phish-food/internal/stocks"

	"log"
)

func archivePost(p []reddit.Post, sub string, conn database.Connection) {
	record := database.NewRedditResponseArchiveTable(p, sub, StartTime)

	err := conn.PutRedditResonseArchiveRecord(record)
	if err != nil {
		log.Printf("Could not archive posts from %s: %s", sub, err.Error())
	}
}

func getPosts(sub string, limit int, auth reddit.AuthToken, conn database.Connection) (p []reddit.Post, err error) {

	p, err = reddit.GetHot(sub, limit, auth)
	if err != nil {
		return p, err
	}

	for i := range p {
		p[i].Comments, err = reddit.FetchPostComments(p[i], auth)
		if err != nil {
			return p, err
		}
	}

	go archivePost(p, sub, conn)

	return p, err
}

func countRef(p []reddit.Post, tickers []stocks.Stock) (r []report.StockReport, err error) {

	r, err = report.CountStockReferences(tickers, p)
	report.SortStockReports(r)

	return r, err
}

package etl

import (
	"github.com/fischersean/phish-food/internal/reddit"
	"github.com/fischersean/phish-food/internal/report"
	"github.com/fischersean/phish-food/internal/stocks"
)

func CountRef(p []reddit.Post, tickers []stocks.Stock) (r []report.StockReport, err error) {

	r, err = report.CountStockReferences(tickers, p)
	report.SortStockReports(r)

	return r, err
}

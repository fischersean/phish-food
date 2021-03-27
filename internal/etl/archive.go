package etl

import (
	"github.com/fischersean/phish-food/internal/reddit"
	"github.com/fischersean/phish-food/internal/report"

	db "github.com/fischersean/phish-food/internal/database"

	"log"
	"time"
)

func ArchivePost(p []reddit.Post, sub string, conn db.Connection) {

	if len(p) == 0 {
		return
	}

	var archiveTime time.Time
	if p[0].DownloadTime.IsZero() {
		archiveTime = time.Now()
	} else {
		archiveTime = p[0].DownloadTime
	}
	record := db.NewRedditResponseArchiveRecord(p, sub, archiveTime)

	err := conn.PutRedditResponseArchiveRecord(record)
	if err != nil {
		log.Printf("Could not archive posts from %s: %s", sub, err.Error())
	}
}

func PutRecord(conn db.Connection, sr []report.StockReport, sub string, t time.Time) (err error) {

	e := db.NewEtlResultsRecord(sr, sub, t)
	err = conn.PutEtlResultsRecord(e)
	return err
}

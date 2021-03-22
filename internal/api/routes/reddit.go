package routes

import (
	"github.com/fischersean/phish-food/internal/api"
	db "github.com/fischersean/phish-food/internal/database"

	"net/http"
	"time"
)

const (
	// Make sure the time is also in GMT
	RawTimeFormat = time.RFC3339
)

func HandleGetLatestRedditData(w http.ResponseWriter, r *http.Request) {

	// Results are delayed by an hour to make sure the ETL pipeline has enough time
	date := time.Now().Add(-time.Hour)

	// Local testing only
	date = date.Add(4 * time.Hour)

	q := r.URL.Query()
	q.Add("datetime", date.Format(RawTimeFormat))

	r.URL.RawQuery = q.Encode()
	HandleGetExactRedditData(w, r)
}

func HandleGetExactRedditData(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	subredditParam := q["subreddit"]
	daterawParam := q["datetime"]
	if len(subredditParam) == 0 || len(daterawParam) == 0 {
		http.Error(w, "Invalid request parameters", 400)
	}

	subreddit := subredditParam[0]
	dateraw := daterawParam[0]
	date, err := time.Parse(RawTimeFormat, dateraw)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	conn := db.SharedConnection
	etlRecordResponse, err := conn.GetEtlResultsRecord(db.EtlResultsQueryInput{
		Subreddit: subreddit,
		Date:      date,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// We are going to do this the slow way to return just the requested hour
	etlRecord := db.EtlResultsRecord{}
	hour := date.Hour()
	for _, v := range etlRecordResponse {
		if v.Hour == hour {
			etlRecord = v
		}
	}

	api.HttpServeMarahallableData(w, etlRecord)
	return
}

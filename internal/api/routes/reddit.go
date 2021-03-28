package routes

import (
	"github.com/fischersean/phish-food/internal/api"
	db "github.com/fischersean/phish-food/internal/database"

	"net/http"
	"time"

	"log"
)

const (
	// Make sure the time is also in GMT
	RawTimeFormat = time.RFC3339
)

func HandleGetLatestRedditData(w http.ResponseWriter, r *http.Request) {

	// Results are delayed by an hour to make sure the ETL pipeline has enough time
	q := r.URL.Query()

	subredditParam := q["subreddit"]
	if len(subredditParam) == 0 {
		http.Error(w, "Invalid request parameters", 400)
	}

	subreddit := subredditParam[0]

	conn := db.SharedConnection
	etlRecord, err := conn.GetLatestEtlResultsRecord(db.EtlResultsQueryInput{
		Subreddit: subreddit,
	})

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not find records", 500)
		return
	}

	_, err = api.HttpServeMarahallableData(w, etlRecord)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
	}
}

func HandleGetExactRedditData(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	subredditParam := q["subreddit"]
	daterawParam := q["datetime"]
	log.Println(daterawParam)
	if len(subredditParam) == 0 || len(daterawParam) == 0 {
		http.Error(w, "Invalid request parameters", 400)
	}

	subreddit := subredditParam[0]
	dateraw := daterawParam[0]
	date, err := time.Parse(RawTimeFormat, dateraw)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

	log.Println(date)

	conn := db.SharedConnection
	etlRecordResponse, err := conn.GetEtlResultsRecord(db.EtlResultsQueryInput{
		Subreddit: subreddit,
		Date:      date,
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Could not find records", 500)
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

	_, err = api.HttpServeMarahallableData(w, etlRecord)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 500)
	}
}

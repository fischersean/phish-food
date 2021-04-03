package routes

import (
	"github.com/fischersean/phish-food/internal/api"
	db "github.com/fischersean/phish-food/internal/database"
	"github.com/fischersean/phish-food/internal/etl"
	"github.com/fischersean/phish-food/internal/router"

	"log"
	"net/http"
	"time"
)

const (
	// Make sure the time is also in GMT
	RawTimeFormat = time.RFC3339
)

// subIsSupported returns whether the subreddit can be found in the ETL database
// Case Sensitive
func subIsSupported(subreddit string) bool {

	var supported bool
	for _, v := range etl.FetchTargets {
		if subreddit == v {
			return true
		}
	}

	return supported
}

func HandleGetLatestRedditData(w http.ResponseWriter, r *http.Request) {

	subreddit := router.GetField(r, 0)
	if subreddit == "" {
		http.Error(w, "Invalid request parameters", 400)
		return
	}
	if !subIsSupported(subreddit) {
		http.Error(w, "Requested subreddit is not supported", 400)
		return
	}

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

	subreddit := router.GetField(r, 0)
	daterawParam := q["datetime"]
	log.Println(daterawParam)
	if subreddit == "" || len(daterawParam) == 0 {
		http.Error(w, "Invalid request parameters", 400)
		return
	}
	if !subIsSupported(subreddit) {
		http.Error(w, "Requested subreddit is not supported", 400)
		return
	}

	dateraw := daterawParam[0]
	date, err := time.Parse(RawTimeFormat, dateraw)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), 400)
		return
	}

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

func HandleGetArchivedRedditData(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	daterawParam := q["datetime"]
	permalinkParam := q["permalink"]
	log.Println(daterawParam)
	if len(permalinkParam) == 0 || len(daterawParam) == 0 {
		http.Error(w, "Invalid request parameters", 400)
		return
	}

	//subreddit := subredditParam[0]
	//if !subIsSupported(subreddit) {
	//http.Error(w, "Requested subreddit is not supported", 400)
	//return
	//}

}

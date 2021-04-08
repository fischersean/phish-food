package routes

import (
	"github.com/fischersean/phish-food/internal/api"
	db "github.com/fischersean/phish-food/internal/database"
	"github.com/fischersean/phish-food/internal/etl"

	"github.com/fischersean/phish-food/internal/router"
	mw "github.com/fischersean/phish-food/internal/router/middleware"

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
	if subreddit == "" || !subIsSupported(subreddit) {
		mw.Error(w, http.StatusBadRequest)
		return
	}

	conn := db.SharedConnection
	etlRecord, err := conn.GetLatestEtlResultsRecord(db.EtlResultsQueryInput{
		Subreddit: subreddit,
	})

	if err != nil {
		log.Println(err.Error())
		mw.Error(w, http.StatusInternalServerError)
		return
	}

	_, err = api.HttpServeMarahallableData(w, etlRecord)
	if err != nil {
		log.Println(err.Error())
	}
}

func HandleGetExactRedditData(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	subreddit := router.GetField(r, 0)
	daterawParam := q["datetime"]
	log.Println(daterawParam)
	if subreddit == "" || len(daterawParam) == 0 || !subIsSupported(subreddit) {
		mw.Error(w, http.StatusBadRequest)
		return
	}

	dateraw := daterawParam[0]
	date, err := time.Parse(RawTimeFormat, dateraw)
	if err != nil {
		// We are returning more useful information on this one since the time format can be confusing
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn := db.SharedConnection
	etlRecordResponse, err := conn.GetEtlResultsRecord(db.EtlResultsQueryInput{
		Subreddit: subreddit,
		Date:      date,
	})
	if err != nil {
		log.Println(err.Error())
		mw.Error(w, http.StatusInternalServerError)
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
	}
}

func HandleGetArchivedRedditData(w http.ResponseWriter, r *http.Request) {

	key := router.GetField(r, 0)
	if key == "" {
		mw.Error(w, http.StatusBadRequest)
		return
	}

	conn := db.SharedConnection
	record, err := conn.GetRedditPostArchiveRecord(db.RedditPostArchiveQueryInput{
		Key: key + ".json",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = api.HttpServeMarahallableData(w, record)
	if err != nil {
		log.Println(err.Error())
	}
}

func HandleListArchivedRedditData(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	daterawParam := q["datetime"]
	subreddit := router.GetField(r, 0)
	if subreddit == "" || len(daterawParam) == 0 || !subIsSupported(subreddit) {
		mw.Error(w, http.StatusBadRequest)
		return
	}

	dateraw := daterawParam[0]
	date, err := time.Parse(RawTimeFormat, dateraw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn := db.SharedConnection
	manifest, err := conn.ListRedditPostArchiveRecord(db.RedditPostArchiveListInput{
		Subreddit: subreddit,
		Date:      date,
	})

	if err != nil {
		log.Println(err.Error())
		mw.Error(w, http.StatusInternalServerError)
		return
	}

	_, err = api.HttpServeMarahallableData(w, manifest)
	if err != nil {
		log.Println(err.Error())
	}

}

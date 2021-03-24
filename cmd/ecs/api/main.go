package main

import (
	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	mw "github.com/fischersean/phish-food/internal/api/middleware"
	"github.com/fischersean/phish-food/internal/api/routes"
	db "github.com/fischersean/phish-food/internal/database"

	_ "github.com/fischersean/phish-food/internal/tzinit"

	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	// AWS Session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	var err error
	db.SharedConnection, err = db.Connect(db.ConnectionInput{
		Session: sess,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	http.Handle("/reddit_latest", mw.ApiKeyRequired(routes.HandleGetLatestRedditData, "/reddit_latest"))
	http.HandleFunc("/", routes.HandleHealthCheck)

	port := os.Getenv("API_PORT")
	if port == "" {
		log.Fatal("Could not find port env variable")
	}

	log.Println("Listening on " + port + "...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}

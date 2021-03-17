package main

import (
	"github.com/aws/aws-sdk-go/aws/session"

	mw "github.com/fischersean/phish-food/internal/api/middleware"
	"github.com/fischersean/phish-food/internal/api/routes"
	db "github.com/fischersean/phish-food/internal/database"

	"log"
	"net/http"
	"os"
)

func main() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	var err error
	db.SharedConnection, err = db.Connect(db.ConnectionInput{
		Session:   sess,
		UserTable: os.Getenv("USER_TABLE"),
	})
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err.Error())
	}

	http.Handle("/app/key", mw.AuthRequired(routes.HandleGetKey))

	log.Fatal(http.ListenAndServe(os.Getenv("API_PORT"), nil))

}

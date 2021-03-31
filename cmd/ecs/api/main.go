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

func apiKeyValidation(r *http.Request) (valid bool, err error) {
	route := r.URL.EscapedPath()
	key := r.Header.Get("x-api-key")
	if key == "" {
		return valid, err
	}

	conn := db.SharedConnection
	keyRecord, err := conn.GetKeyPermissions(db.ApiKeyQueryInput{
		UnhashedKey: key,
	})
	if err != nil || !keyRecord.Enabled {
		log.Printf("Authentication Error: %s", err.Error())
		return valid, err
	}

	for _, v := range keyRecord.Permissions {
		if v == route {
			valid = true
			break
		}
	}

	return valid, err
}

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

	routeOptions := mw.HandlerOptions{
		Methods: []string{
			http.MethodGet,
		},
		Cors: mw.CorsOptions{
			Enabled: true,
			AllowHeaders: []string{
				"Access-Control-Allow-Origin",
				"Content-Type",
				"x-api-key",
			},
			AllowOrigin: []string{
				"*",
			},
		},
		Authentication: mw.AuthenticationOptions{
			Required:       true,
			ValidationFunc: apiKeyValidation,
		},
	}
	http.Handle("/reddit", mw.Register(routes.HandleGetExactRedditData, routeOptions))
	http.Handle("/reddit_latest", mw.Register(routes.HandleGetLatestRedditData, routeOptions))

	http.HandleFunc("/", routes.HandleHealthCheck)

	port := os.Getenv("API_PORT")
	if port == "" {
		log.Fatal("Could not find port env variable")
	}

	log.Println("Listening on " + port + "...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))

}

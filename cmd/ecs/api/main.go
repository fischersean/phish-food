package main

import (
	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/fischersean/phish-food/internal/api/routes"
	db "github.com/fischersean/phish-food/internal/database"
	"github.com/fischersean/phish-food/internal/router"
	mw "github.com/fischersean/phish-food/internal/router/middleware"

	_ "github.com/fischersean/phish-food/internal/tzinit"

	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
)

func apiKeyValidation(r *http.Request) (valid bool) {

	route := r.URL.EscapedPath()
	key := r.Header.Get("x-api-key")
	if key == "" {
		return valid
	}

	conn := db.SharedConnection
	keyRecord, err := conn.GetKeyPermissions(db.ApiKeyQueryInput{
		UnhashedKey: key,
	})
	if err != nil || !keyRecord.Enabled {
		log.Printf("Authentication Error: %s", err.Error())
		return valid
	}

	for _, v := range keyRecord.Permissions {
		r, err := regexp.Compile(v)
		if err != nil {
			log.Printf("Could not compile permissions regex: %s", v)
			return
		}
		if len(r.FindString(route)) > 0 {
			valid = true
			break
		}
	}

	return valid
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

	router.Handle("/reddit/*([^/]+)/latest", mw.WithOptions(routes.HandleGetLatestRedditData, routeOptions))
	router.Handle("/reddit/*([^/]+)", mw.WithOptions(routes.HandleGetExactRedditData, routeOptions))

	router.Handle("/reddit/archive/*([a-zA-Z_0-9/]+).json", mw.WithOptions(routes.HandleGetArchivedRedditData, routeOptions))
	router.Handle("/reddit/archive/*([^/]+)/manifest", mw.WithOptions(routes.HandleListArchivedRedditData, routeOptions))

	router.HandleFunc("/", routes.HandleHealthCheck)

	port := os.Getenv("API_PORT")
	if port == "" {
		log.Fatal("Could not find port env variable")
	}

	log.Println("Listening on " + port + "...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), http.HandlerFunc(router.Serve)))

}

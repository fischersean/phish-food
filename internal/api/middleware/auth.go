package middleware

import (
	db "github.com/fischersean/phish-food/internal/database"
	"net/http"
)

// ApiKeyRequired validates that the supplied x-api-key header gives the requestor access to the route
func ApiKeyRequired(h func(w http.ResponseWriter, r *http.Request), route string) http.Handler {
	return apiKeyRequiredHandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r) // call original
	}), route)
}

func apiKeyRequiredHandlerFunc(h http.HandlerFunc, route string) http.Handler {
	return withCORSHandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := r.Header.Get("x-api-key")
		if key == "" {
			http.Error(w, "Permission Denied", 401)
			return
		}

		conn := db.SharedConnection
		keyRecord, err := conn.GetKeyPermissions(db.ApiKeyQueryInput{
			UnhashedKey: key,
		})
		if err != nil || !keyRecord.Enabled {
			// Although this may be an internal server error, we will not reveal that
			// incase someone is tryig to brute force a key
			http.Error(w, "Permission Denied", 401)
			return
		}

		valid := false
		for _, v := range keyRecord.Permissions {
			if v == route {
				valid = true
				break
			}
		}

		if !valid {
			http.Error(w, "Permission Denied", 401)
			return
		}

		h(w, r) // call original
	}))
}

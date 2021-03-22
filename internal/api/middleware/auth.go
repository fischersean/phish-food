package middleware

import (
	"net/http"
)

func ApiKeyRequired(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return apiKeyRequiredHandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r) // call original
	}))
}

func apiKeyRequiredHandlerFunc(h http.HandlerFunc) http.Handler {
	return withCORSHandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Check to see if a valid api key Authorization header is present
		h(w, r) // call original
	}))
}

package middleware

import (
	"net/http"
	"os"
)

func WithCORS(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return withCORSHandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r) // call original
	}))
}

func withCORSHandlerFunc(h http.HandlerFunc) http.Handler {
	return loggedHandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("APP_URL"))
		w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Add("Access-Control-Allow-Headers", "Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		h(w, r) // call original
	}))
}

func AuthRequired(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return withCORSHandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: This function should also validate the JWT is the correct length/general format
		tokenHeader := r.Header["Authorization"]
		if len(tokenHeader) != 1 {
			http.Error(w, "Invalid token header", 400)
			return
		}
		h(w, r)
	}))
}

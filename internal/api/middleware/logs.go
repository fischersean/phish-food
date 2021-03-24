package middleware

import (
	"log"
	"net/http"
)

func logRequest(r *http.Request) {
	log.Printf("URL: %s | Method: %s | Host: %s", r.URL.EscapedPath(), r.Method, r.Host)
}

//func logResponse(w http.ResponseWriter) {
//log.Printf("Resonse Headers: %s", w.Header())
//}

// Logged logs the details of the request and response
func Logged(h func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return loggedHandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r) // call original
	}))
}

// LoggedNaked logs the details of the request and response
func loggedHandlerFunc(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		h(w, r) // call original
		//logResponse(w)
	})
}

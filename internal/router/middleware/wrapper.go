package middleware

import (
	"net/http"
)

// Error is a shorthand for writing standard error messages to the response writer
func Error(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func WithOptions(h http.HandlerFunc, options HandlerOptions) http.Handler {

	// Options are added in the reverse order that they need to be evaluated at request time
	// Authentication always happens after CORS
	if options.Authentication.Required {
		h = withAuthentication(h, options.Authentication)
	}

	// We will always validate the request method
	h = withMethodValidation(h, options.Methods)

	// CORS needs to happen right after the request is logged
	if options.Cors.Enabled {
		h = withCors(h, options.Cors)
	}

	// Requests are always logged regardless of any other options
	return withLogging(h)
}

func withMethodValidation(h http.HandlerFunc, methods []string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var methodFound bool
		for _, m := range methods {
			if m == r.Method {
				methodFound = true
			}
		}
		if methodFound {
			h(w, r) // call original
		} else {
			Error(w, http.StatusMethodNotAllowed)
		}
	})
}

func withLogging(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		h(w, r) // call original
	})
}

func withAuthentication(h http.HandlerFunc, options AuthenticationOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validated := options.ValidationFunc(r)
		if validated {
			h(w, r)
		} else {
			Error(w, http.StatusForbidden)
		}
	})
}

func withCors(h http.HandlerFunc, options CorsOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var requestOriginAllowed bool
		var requestMethodAllowed bool

		for _, origin := range options.AllowOrigins {
			w.Header().Add("Access-Control-Allow-Origin", origin)
			if origin == "*" || origin == r.Host {
				requestOriginAllowed = true
			}
		}

		for _, header := range options.AllowHeaders {
			w.Header().Add("Access-Control-Allow-Headers", header)
		}

		// Implicitly add OPTIONS method
		w.Header().Add("Access-Control-Allow-Methods", http.MethodOptions)
		for _, method := range options.AllowMethods {
			w.Header().Add("Access-Control-Allow-Methods", method)
			if method == r.Method {
				requestMethodAllowed = true
			}
		}

		if r.Method == http.MethodOptions {
			return
		}

		if !requestOriginAllowed {
			Error(w, http.StatusForbidden)
			return
		}
		if !requestMethodAllowed {
			Error(w, http.StatusMethodNotAllowed)
			return
		}

		h(w, r) // call original
	})
}

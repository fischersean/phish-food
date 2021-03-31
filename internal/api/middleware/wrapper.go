package middleware

import (
	"net/http"
)

func Register(h http.HandlerFunc, options HandlerOptions) http.Handler {

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
	return http.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var methodFound bool
		for _, m := range methods {
			if m == r.Method {
				methodFound = true
			}
		}
		if methodFound {
			h(w, r) // call original
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))
}

func withLogging(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		h(w, r) // call original
	}))
}

func withAuthentication(h http.HandlerFunc, options AuthenticationOptions) http.HandlerFunc {
	return http.HandlerFunc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		validated, err := options.ValidationFunc(r)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if validated {
			h(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	}))
}

func withCors(h http.HandlerFunc, options CorsOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		for _, origin := range options.AllowOrigin {
			w.Header().Add("Access-Control-Allow-Origin", origin)
		}

		for _, header := range options.AllowHeaders {
			w.Header().Add("Access-Control-Allow-Headers", header)
		}

		if r.Method == http.MethodOptions {
			return
		}

		h(w, r) // call original
	})
}

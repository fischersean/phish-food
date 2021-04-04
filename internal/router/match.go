package router

import (
	"context"
	"net/http"
	"regexp"
)

type route struct {
	regex   *regexp.Regexp
	handler http.Handler
}

var routes []route

func Handle(pattern string, h http.Handler) {
	r := regexp.MustCompile("^" + pattern + "$")
	routes = append(routes, route{
		regex:   r,
		handler: h,
	})
}

func HandleFunc(pattern string, h http.HandlerFunc) {
	Handle(pattern, http.Handler(h))
}

func Serve(w http.ResponseWriter, r *http.Request) {
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			// 1: because we dont care about the whole route, only the variable parts
			ctx := context.WithValue(r.Context(), CtxKey{}, matches[1:])
			route.handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	}
	http.NotFound(w, r)
}

type CtxKey struct{}

func GetField(r *http.Request, index int) string {
	fields := r.Context().Value(CtxKey{}).([]string)
	if len(fields) < index-1 {
		return ""
	}
	return fields[index]
}

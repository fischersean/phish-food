package middleware

import (
	"net/http"
)

type HandlerOptions struct {
	Methods        []string
	Cors           CorsOptions
	Authentication AuthenticationOptions
}

type CorsOptions struct {
	Enabled      bool
	AllowHeaders []string
	AllowOrigins []string
	AllowMethods []string
}

type AuthenticationOptions struct {
	Required bool
	// ValidationFunc should return true or false based on whether the request is authenticated
	ValidationFunc func(*http.Request) bool
}

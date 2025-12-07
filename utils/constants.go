package utils

import "net/http"

const (
	APP_NAME            = "Curl Wrapper"
	DEFAULT_URL_EXAMPLE = "http://localhost:5172"
)

type HTTPMethod uint8

const (
	METHOD_GET HTTPMethod = iota
	METHOD_POST
	METHOD_DELETE
	METHOD_PATCH
	METHOD_PUT
	METHOD_OPTIONS
)

var (
	ALL_HTTP_METHODS = []string{"GET", "POST", "DELETE", "PATCH", "PUT", "OPTIONS"}

	HTTPMethodMap = map[HTTPMethod]string{
		METHOD_GET:     http.MethodGet,
		METHOD_POST:    http.MethodPost,
		METHOD_DELETE:  http.MethodDelete,
		METHOD_PATCH:   http.MethodPatch,
		METHOD_PUT:     http.MethodPut,
		METHOD_OPTIONS: http.MethodOptions,
	}
)

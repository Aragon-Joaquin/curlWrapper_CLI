package main

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"time"
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

	httpMethodMap = map[HTTPMethod]string{
		METHOD_GET:     http.MethodGet,
		METHOD_POST:    http.MethodPost,
		METHOD_DELETE:  http.MethodDelete,
		METHOD_PATCH:   http.MethodPatch,
		METHOD_PUT:     http.MethodPut,
		METHOD_OPTIONS: http.MethodOptions,
	}

	customClient = &http.Client{
		Timeout: time.Second * 4,
	}

	errorWrongMethod = errors.New("unknown method selected")
	errorBadURL      = errors.New("the url provided is wrong")
)

func MakeHTTPCall(URL string, httpMethod HTTPMethod) error {
	goodUrl, err := url.Parse(URL)

	if err != nil {
		slog.Warn("Error while parsing URL", "Err", err.Error())
		return errorBadURL
	}

	method, ok := httpMethodMap[httpMethod]

	if !ok {
		slog.Warn("Error while parsing METHOD", "Method", httpMethod)
		return errorWrongMethod
	}

	resp, err := customClient.Do(&http.Request{
		Method: method,
		URL:    goodUrl,
	})

	if err != nil {
		slog.Warn("Error while doing request", "Err", err.Error())
		return err
	}

	defer resp.Body.Close()

	return nil
}

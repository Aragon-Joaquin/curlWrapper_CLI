package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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

type RequestJson struct {
	*http.Response
	ParsedBody string
}

func MakeHTTPCall() (*RequestJson, error) {
	goodUrl, err := url.Parse(GlobalFieldState.URLField)

	if err != nil {
		slog.Error("Error while parsing URL", "Err", err.Error())
		return nil, errorBadURL
	}

	method, ok := httpMethodMap[GlobalFieldState.MethodField]

	if !ok {
		slog.Error("Error while parsing METHOD", "Method", GlobalFieldState.MethodField)
		return nil, errorWrongMethod
	}

	resp, err := customClient.Do(&http.Request{
		Method: method,
		URL:    goodUrl,
	})

	if err != nil {
		slog.Error("Error while doing request", "Err", err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		slog.Error("Error while reading body", "Err", err.Error())
		return nil, err
	}

	var prettyJSON bytes.Buffer

	if IsJSON(string(body)) {
		err = json.Indent(&prettyJSON, body, "", "\t")

		if err != nil {
			slog.Error("Error while prettifying JSON", "Err", err.Error(), "JSON", body)
			return nil, err
		}

		return &RequestJson{
			ParsedBody: prettyJSON.String(),
			Response:   resp,
		}, nil
	}

	return &RequestJson{
		ParsedBody: string(body),
		Response:   resp,
	}, nil
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

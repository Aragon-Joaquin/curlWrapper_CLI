package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	ut "github.com/Aragon-Joaquin/curlWrapper_CLI/utils"
)

var (
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

	method, ok := ut.HTTPMethodMap[GlobalFieldState.MethodField]

	if !ok {
		slog.Error("Error while parsing METHOD", "Method", GlobalFieldState.MethodField)
		return nil, errorWrongMethod
	}

	httpRequest := &http.Request{
		Method: method,
		URL:    goodUrl,
	}

	if method != http.MethodGet && GlobalFieldState.Body != "" {
		httpRequest.Header = http.Header{
			"Content-Type": []string{"application/json", "charset=utf-8"},
		}

		r := io.NopCloser(strings.NewReader(GlobalFieldState.Body))
		defer r.Close()

		httpRequest.Body = r
	}

	resp, err := customClient.Do(httpRequest)

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

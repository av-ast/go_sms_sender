package main

import (
	"gopkg.in/resty.v1"
	"net/http"
)

func send_request(gateway_params map[string]string, query_params map[string]string) *resty.Response {
	if gateway_params["enabled"] == "true" {
		resp, err := resty.R().SetQueryParams(query_params).Get(gateway_params["url"])

		if err != nil {
			panic(err)
		}

		return resp
	} else {
		raw_resp := &http.Response{Status: "200 OK", StatusCode: http.StatusOK, Body: http.NoBody}
		resp := &resty.Response{RawResponse: raw_resp}

		return resp
	}
}

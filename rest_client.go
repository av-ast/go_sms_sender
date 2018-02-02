package main

import (
	"gopkg.in/resty.v1"
	"net/http"
)

func send_request(gateway_params map[string]string, query_params map[string]string) *resty.Response {
	var resp *resty.Response
	var err interface{}

	if gateway_params["enabled"] == "true" {
		resp, err = resty.R().SetQueryParams(query_params).Get(gateway_params["url"])

		if err != nil {
			panic(err)
		}
	} else {
		raw_resp := new(http.Response)
		raw_resp.Status = "200 OK"
		raw_resp.StatusCode = http.StatusOK
		raw_resp.Body = http.NoBody

		resp = new(resty.Response)
		resp.RawResponse = raw_resp
	}

	return resp
}

package main

import (
	"gopkg.in/resty.v1"
)

func send_request(url string, query_params map[string]string) *resty.Response {
	resp, err := resty.R().SetQueryParams(query_params).Get(url)

	if err != nil {
		panic(err)
	}

	return resp
}

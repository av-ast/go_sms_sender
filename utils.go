package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type JsonData map[string]string

func jsonFromRequestBody(request *http.Request) JsonData {
	var data JsonData

	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	return data
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

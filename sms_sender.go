package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type JsonData map[string]interface{}

type Response struct {
	Status string   `json:"status"`
	Data   JsonData `json:"data"`
}

type JsonResponse interface {
	toJson() string
}

var gateway_params map[string]string

func main() {
	addr := getenv("LISTEN_ADDR", ":8000")
	gateway_params = make(map[string]string)
	gateway_params["url"] = getenv("SMS_GATEWAY_URL", "https://bsms.tele2.ru/api")
	gateway_params["login"] = getenv("SMS_GATEWAY_LOGIN", "")
	gateway_params["password"] = getenv("SMS_GATEWAY_PASSWORD", "")

	router := httprouter.New()
	router.POST("/sms", send_sms)
	// router.GET("/sms/:id", get_sms_status)

	log.Print("Starting SMS sender on ", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Could not run SMS sender: ", err)
	}
}

func send_sms(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	defer request.Body.Close()
	log.Print("Request params: ", params)

	data := jsonFromRequestBody(request)
	log.Print("Request body: ", data)

	url := fmt.Sprintf("%s?operation=send&login=%s&password=%s&shortcode=%s&msisdn=7%s&text=%s",
		gateway_params["url"],
		gateway_params["login"],
		gateway_params["password"],
		data["from"],
		data["to"],
		data["text"])
	log.Print("URL: ", url)

	gw_resp, err := http.Get(url)
	log.Print("Gateway Response: ", gw_resp)

	if err != nil {
		panic(err)
	}

	defer gw_resp.Body.Close()
	body, err := ioutil.ReadAll(gw_resp.Body)

	if err != nil {
		panic(err)
	}

	log.Print("Gateway Response body: ", fmt.Sprintf("%s", body))

	resp := &Response{Status: "success", Data: data}
	fmt.Fprintf(writer, resp.toJson())
}

func (resp Response) toJson() string {
	result, err := json.Marshal(resp)

	if err != nil {
		panic(err)
	}

	return string(result)
}

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

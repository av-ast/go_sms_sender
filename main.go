package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"regexp"
)

type Response struct {
	Status string   `json:"status"`
	Data   JsonData `json:"data"`
}

type JsonResponse interface {
	toJson() string
}

const sms_id_format = "^[0-9]+$"

var (
	gateway_params map[string]string
	sms_id_regexp  *regexp.Regexp
)

func main() {
	addr := getenv("LISTEN_ADDR", ":8000")
	gateway_params = make(map[string]string)
	gateway_params["enabled"] = getenv("SMS_GATEWAY_ENABLED", "false")
	gateway_params["url"] = getenv("SMS_GATEWAY_URL", "https://bsms.tele2.ru/api")
	gateway_params["login"] = getenv("SMS_GATEWAY_LOGIN", "")
	gateway_params["password"] = getenv("SMS_GATEWAY_PASSWORD", "")
	gateway_params["shortcode"] = getenv("SMS_GATEWAY_SHORTCODE", "")
	sms_id_regexp, _ = regexp.Compile(sms_id_format)

	router := httprouter.New()
	router.POST("/sms", send_sms)
	router.GET("/sms/:id", get_sms_status)

	log.Print("Starting SMS sender on ", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("Could not run SMS sender: ", err)
	}
}

func send_sms(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	defer request.Body.Close()

	sms_data := jsonFromRequestBody(request)
	query_params := map[string]string{
		"operation": "send",
		"login":     gateway_params["login"],
		"password":  gateway_params["password"],
		"shortcode": gateway_params["shortcode"],
		"msisdn":    fmt.Sprintf("7%s", sms_data["phone"]),
		"text":      sms_data["text"],
	}

	log.Print("Request url: ", gateway_params["url"])
	log.Print("Request params: ", query_params)

	gw_resp := send_request(gateway_params, query_params)

	log.Print("Gateway Response Status: ", gw_resp.Status())
	log.Print("Gateway Response Body: ", gw_resp.String())

	resp := &Response{Data: make(map[string]string)}

	if gw_resp.StatusCode() == http.StatusOK && sms_id_regexp.MatchString(gw_resp.String()) {
		resp.Status = "success"
		resp.Data["sms_id"] = gw_resp.String()
	} else {
		resp.Status = "error"
		resp.Data["error"] = gw_resp.String()
	}

	fmt.Fprintf(writer, resp.toJson())
}

func get_sms_status(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	defer request.Body.Close()

	query_params := map[string]string{
		"operation": "status",
		"login":     gateway_params["login"],
		"password":  gateway_params["password"],
		"id":        params.ByName("id"),
	}

	log.Print("Request url: ", gateway_params["url"])
	log.Print("Request params: ", query_params)

	gw_resp := send_request(gateway_params, query_params)

	log.Print("Gateway Response Status: ", gw_resp.Status())
	log.Print("Gateway Response Body: ", gw_resp.String())

	resp := &Response{Data: make(map[string]string)}

	if gw_resp.StatusCode() == http.StatusOK {
		resp.Status = "success"
		resp.Data["delivery_status"] = gw_resp.String()
	} else {
		resp.Status = "error"
		resp.Data["error"] = gw_resp.String()
	}

	fmt.Fprintf(writer, resp.toJson())
}

func (resp Response) toJson() string {
	result, err := json.Marshal(resp)

	if err != nil {
		panic(err)
	}

	return string(result)
}

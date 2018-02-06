package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"regexp"
)

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

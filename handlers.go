package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type Response struct {
	Status string   `json:"status"`
	Data   JsonData `json:"data"`
}

type JsonResponse interface {
	toJson() string
}

/**
 * @api {post} /sms Send SMS
 * @apiName SendSms
 * @apiGroup Sms
 * @apiHeader {string=application/json} Content-Type
 *
 * @apiParam {String} phone Phone number.
 * @apiParam {String} text SMS text.
 *
 * @apiSuccessExample {json} 200 Success
 * {
 *   "data": {"sms_id": "3243242423"},
 *   "status": "success"
 * }
 */
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

/**
* @api {get} /sms/:id Get SMS status
 * @apiName GetSmsStatus
 * @apiGroup Sms
 *
 * @apiParam {String} id SMS id.
 *
 * @apiSuccessExample {json} 200 Success
 * {
 *   "data": {"delivery_status": "delivered"},
 *   "status": "success"
 * }
*/
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

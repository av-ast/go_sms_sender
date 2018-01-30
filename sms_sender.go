package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

type ResponseData map[string]interface{}

type Response struct {
	Status string       `json:"status"`
	Data   ResponseData `json:"data"`
}

type JsonResponse interface {
	toJson() string
}

var (
	storage ResponseData
	addr    = getenv("LISTEN_ADDR", ":8000")
)

func main() {
	storage = ResponseData{}

	router := httprouter.New()
	router.GET("/entries", show)
	router.GET("/entries/:key", show)
	router.PUT("/entries/:key/:value", update)

	log.Print("Starting SMS sender on ", addr)
	err := http.ListenAndServe(addr, router)

	if err != nil {
		log.Fatal("Could not run SMS sender: ", err)
	}
}

func show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	k := p.ByName("key")
	resp := &Response{Status: "success"}

	if k == "" {
		log.Print("Show collection: ", p, storage)
		resp.Data = storage
	} else {
		log.Print("Show entry: ", p, storage[k])
		data := make(ResponseData)
		data[k] = storage[k]
		resp.Data = data
	}

	fmt.Fprintf(w, resp.toJson())
}

func update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	k := p.ByName("key")
	v := p.ByName("value")
	resp := &Response{Status: "success"}

	storage[k] = v

	log.Print("Update collection: ", p, storage)
	data := make(ResponseData)
	data["key"], data["value"] = k, v

	resp.Data = data
	fmt.Fprintf(w, resp.toJson())
}

func (resp Response) toJson() string {
	result, err := json.Marshal(resp)

	if err != nil {
		panic(err)
	}

	return string(result)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

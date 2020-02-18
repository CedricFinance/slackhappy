package main

import (
	"github.com/CedricFinance/slackhappy"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		payload, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Panic(err)
		}
		err = slackhappy.OnPubSubMessage(request.Context(), slackhappy.PubSubMessage{
			Data: payload,
		})
		log.Printf("failed to process the PubSub message: %q", err)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:8889", nil))
}

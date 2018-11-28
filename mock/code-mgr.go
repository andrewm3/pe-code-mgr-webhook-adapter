package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andrewm3/pe-code-mgr-webhook-adapter/adapter"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		event, err := adapter.ParseEvent(r)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		fmt.Printf("Push to %s received with parameters: %s\n", event.Ref, r.URL.Query())
	})

	log.Fatal(http.ListenAndServe(":8170", nil))
}

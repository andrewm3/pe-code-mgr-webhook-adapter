package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/andrewm3/pe-code-mgr-webhook-adapter/adapter"
)

var port int
var redirect string
var whitelist = []string{
	"refs/heads/production",
}

func init() {
	flag.IntVar(&port, "port", 8080, "The port to serve on")
	flag.StringVar(
		&redirect,
		"redirect",
		"https://localhost:8170/code-manager/v1/webhook",
		"The URL for Code Manager, which requests will be forwarded to",
	)
}

func main() {
	flag.Parse()
	config := adapter.HandlerConfig{
		Redirect:  redirect,
		Whitelist: whitelist,
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		adapter.EventHandler(w, r, config)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

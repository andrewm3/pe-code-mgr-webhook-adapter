package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/andrewm3/pe-code-mgr-webhook-adapter/adapter"
)

var port int
var codeMgrURL string
var whitelist string

func init() {
	flag.IntVar(&port, "port", 8080, "The port to serve on")
	flag.StringVar(
		&codeMgrURL,
		"code_mgr_url",
		"https://localhost:8170/code-manager/v1/webhook",
		"The URL for Code Manager, which requests will be forwarded to",
	)
	flag.StringVar(&whitelist, "whitelist", "production", "A comma-separated list of allowed branches")
}

func main() {
	flag.Parse()
	config := adapter.HandlerConfig{
		CodeMgrURL: codeMgrURL,
		Whitelist:  strings.Split(whitelist, ","),
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		adapter.EventHandler(w, r, config)
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

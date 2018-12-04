package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// HandlerConfig holds all the required config for the EventHandler
type HandlerConfig struct {
	CodeMgrURL string
	Whitelist  []string
}

// Event is a agnostic representation of the important fields delivered by the webhook
type Event struct {
	Ref string
}

// EventHandler handles all requests
func EventHandler(w http.ResponseWriter, r *http.Request, config HandlerConfig) {
	var allowed bool
	var err error

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	endpoint := fmt.Sprintf("%s%s", config.CodeMgrURL, r.URL)
	u, err := url.Parse(endpoint)
	if err != nil {
		http.Error(w, fmt.Sprintf("Provided URL '%s' cannot be parsed", endpoint), 400)
		return
	}

	replay := CreateReplay(r)
	replay.URL = u

	event, err := ParseEvent(r)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	for _, x := range config.Whitelist {
		if event.Ref == fmt.Sprintf("refs/heads/%s", x) {
			allowed = true
			break
		}
	}

	if allowed {
		var netClient = &http.Client{Timeout: time.Second * 10}
		_, err = netClient.Do(&replay)
		if err != nil {
			http.Error(w, "Error replaying request", 400)
			return
		}
		fmt.Printf("Push to '%s' forwarded\n", event.Ref)
	} else {
		fmt.Printf("Push to '%s' ignored\n", event.Ref)
	}
}

// CreateReplay duplicates a request allowing it to be replayed
func CreateReplay(r *http.Request) http.Request {
	// Copy the original request so it can be replayed
	replay := *r
	replay.RequestURI = ""

	// Artificially duplicate the reader stream in Body
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	r.Body = ioutil.NopCloser(tee)
	replay.Body = ioutil.NopCloser(&buf)

	return replay
}

// ParseEvent takes a request and parses it into an Event
func ParseEvent(r *http.Request) (Event, error) {
	eventType := r.URL.Query().Get("type")
	switch eventType {
	case "github", "gitlab":
		return ParseGenericEvent(r)
	}
	return Event{}, fmt.Errorf("Unsupported type '%s' provided", eventType)
}

// ParseGenericEvent parses a request in the 'default' format
func ParseGenericEvent(r *http.Request) (Event, error) {
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return Event{}, err
	}
	return event, nil
}

package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// HandlerConfig holds all the required config for the EventHandler
type HandlerConfig struct {
	Redirect  string
	Whitelist []string
}

// Event is a agnostic representation of the important fields delivered by the webhook
type Event struct {
	Ref string
}

// GithubEvent holds the fields unmarshaled from a Github webhook payload
type GithubEvent struct {
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

	endpoint := fmt.Sprintf("%s%s", config.Redirect, r.URL)
	u, err := url.Parse(endpoint)
	if err != nil {
		http.Error(w, fmt.Sprintf("Provided URL '%s' cannot be parsed", endpoint), 400)
		return
	}

	// Copy the original request so it can be replayed
	replayed := *r
	replayed.RequestURI = ""
	replayed.URL = u

	// Artificially duplicate the reader stream in Body
	body, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	replayed.Body = ioutil.NopCloser(bytes.NewReader(body))

	event, err := ParseEvent(r)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	for _, x := range config.Whitelist {
		if event.Ref == x {
			allowed = true
			break
		}
	}

	if allowed {
		var netClient = &http.Client{Timeout: time.Second * 10}
		_, err = netClient.Do(&replayed)
		if err != nil {
			http.Error(w, "Error replaying request", 400)
			return
		}
		fmt.Printf("Push to '%s' forwarded\n", event.Ref)
	} else {
		fmt.Printf("Push to '%s' ignored\n", event.Ref)
	}
}

// ParseEvent takes a request and parses it into an Event
func ParseEvent(r *http.Request) (Event, error) {
	eventType := r.URL.Query().Get("type")
	switch eventType {
	case "github":
		return ParseGithubEvent(r)
	}
	return Event{}, fmt.Errorf("Unsupported type '%s' provided", eventType)
}

// ParseGithubEvent parses a request in the expected format from Github
func ParseGithubEvent(r *http.Request) (Event, error) {
	var event GithubEvent
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return Event{}, err
	}
	return Event{event.Ref}, nil
}

package adapter

import (
	"net/http"
	"os"
	"testing"
)

func TestParseGithubEvent(t *testing.T) {
	file, err := os.Open("../test/github.json")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "", file)
	if err != nil {
		t.Fatal(err)
	}

	expected := Event{"refs/tags/simple-tag"}
	actual, err := ParseGithubEvent(req)
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected {
		t.Errorf("ParseGithubEvent returned unexpected value: got %v want %v", actual, expected)
	}
}

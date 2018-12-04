package adapter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestCreateReplay(t *testing.T) {
	file, err := os.Open("../test/github.json")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "", file)
	if err != nil {
		t.Fatal(err)
	}

	replay := CreateReplay(req)
	actual, _ := ioutil.ReadAll(req.Body)
	expected, err := ioutil.ReadAll(replay.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(actual, expected) {
		t.Error("Replayed request does not have the same body")
	}
}

func BenchmarkCreateReplay(b *testing.B) {
	file, _ := os.Open("../test/github.json")
	req, _ := http.NewRequest("POST", "", file)

	for n := 0; n < b.N; n++ {
		CreateReplay(req)
	}
}

func TestParseEvent(t *testing.T) {
	var tests = []struct {
		file     string
		postType string
		expected string
	}{
		{"../test/github.json", "github", "refs/tags/simple-tag"},
		{"../test/gitlab.json", "gitlab", "refs/heads/production"},
	}

	for _, test := range tests {
		file, err := os.Open(test.file)
		if err != nil {
			t.Fatal(err)
		}

		endpoint := fmt.Sprintf("https://localhost:8170/code-manager/v1/webhook/?type=%s", test.postType)
		req, err := http.NewRequest("POST", endpoint, file)
		if err != nil {
			t.Fatal(err)
		}

		expected := Event{test.expected}
		actual, err := ParseEvent(req)
		if err != nil {
			t.Fatal(err)
		}

		if actual != expected {
			t.Errorf("ParseEvent returned unexpected value: got %v want %v", actual, expected)
		}
	}
}

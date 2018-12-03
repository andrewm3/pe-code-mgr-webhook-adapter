package adapter

import (
	"bytes"
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

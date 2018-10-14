package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
)

func TestInsertARoot(t *testing.T) {
	params := url.Values{"id": {uuid.New().String()}, "parent": {"123e4567-e89b-12d3-a456-426655440003"}}
	resp, err := http.PostForm("http://localhost:9000/hierarchy", params)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(body))
}

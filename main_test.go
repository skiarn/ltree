package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
)

//TODO:
//Make 1 million rootnodes, with 10 children and their children have 20 children and thier chidren have 30 childen.
func TestInsertARoot(t *testing.T) {
	params := url.Values{"id": {uuid.New().String()}, "parent": {"f6db4fb7-9148-4528-b776-a00624161de3"}}
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

func TestInsertBig(t *testing.T) {
	parent := "58566036-26d7-4db9-91d6-5afe94583d62"
	for i := 0; i < 10; i++ {
		child1 := doRequest(parent, t)
		for l1 := 0; l1 < 10; l1++ {
			child2 := doRequest(child1, t)
			for l2 := 0; l2 < 10; l2++ {
				child3 := doRequest(child2, t)
				for l3 := 0; l3 < 10; l3++ {
					_ = doRequest(child3, t)
				}
			}
		}
	}
}

func doRequest(parent string, t *testing.T) string {
	params := url.Values{"id": {uuid.New().String()}, "parent": {parent}}
	resp, err := http.PostForm("http://localhost:9000/hierarchy", params)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

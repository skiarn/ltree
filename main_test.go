package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/google/uuid"
)

//TODO:
//Make 1 million rootnodes, with 10 children and their children have 20 children and thier chidren have 30 childen.
func TestInsertARoot(t *testing.T) {
	params := url.Values{"id": {uuid.New().String()}}
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
	makeNodesForARootNode(parent, t)
}

func TestInsertOneMillion(t *testing.T) {

	total := 1000000
	concurrency := 50
	requests := 0
	var wg sync.WaitGroup
	for requests < total {
		if requests+concurrency <= total {
			concurrency = total - requests
		}
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func() {
				rootNode := doRequest("", t)
				makeNodesForARootNode(rootNode, t)
				requests++
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func makeNodesForARootNode(rootnode string, t *testing.T) {
	for i := 0; i < 10; i++ {
		child1 := doRequest(rootnode, t)
		for l1 := 0; l1 < 20; l1++ {
			child2 := doRequest(child1, t)
			for l2 := 0; l2 < 30; l2++ {
				_ = doRequest(child2, t)
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

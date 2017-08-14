package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err.Error())
	}

	text = strings.Trim(text, "\n")
	if text != "" {
		if text == "git branch" {
			listBranches()
		} else if text == "git pr" {
			pullRequest()
		}
	}
}

type Branch struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
		URL string `json:"url"`
	} `json:"commit"`
	Protected     bool   `json:"protected"`
	ProtectionURL string `json:"protection_url"`
}

func listBranches() {
	api := "https://api.github.com/repos/reinarduswindy/gh/branches"

	resp, err := http.Get(api)
	if err != nil {
		panic(err.Error())
	}

	res, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err.Error())
	}

	var listBranches []Branch
	err = json.Unmarshal(res, &listBranches)
	if err != nil {
		panic(err.Error())
	}

	for _, b := range listBranches {
		fmt.Println(b.Name)
	}
}

type PullRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

func pullRequest() {
	api := "https://api.github.com/repos/reinarduswindy/gh/pulls"

	pr := PullRequest{
		Title: "Hello, Pull Request!",
		Body:  "Please pull this in",
		Head:  "reinarduswindy:test",
		Base:  "master",
	}

	jsonData, err := json.Marshal(pr)
	if err != nil {
		panic(err.Error())
	}

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(body))
}

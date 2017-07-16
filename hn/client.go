package hn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	topStoriesAPI = "https://hacker-news.firebaseio.com/v0/topstories.json"
	itemAPI       = "https://hacker-news.firebaseio.com/v0/item/%d.json"
)

type Client interface {
	TopStories() ([]int, error)
	Story(id int) (*Story, error)
}

type Story struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	URL     string `json:"url"`
	Author  string `json:"by"`
	Score   int    `json:"score"`
	Created int64  `json:"time"`
}

func NewClient(h *http.Client) Client {
	return &baseClient{h: h}
}

type baseClient struct {
	h *http.Client
}

func (c *baseClient) TopStories() ([]int, error) {
	rsp, err := c.h.Get(topStoriesAPI)

	if rsp != nil {
		defer rsp.Body.Close()
	}

	if err != nil {
		return []int{}, fmt.Errorf("Client.TopStories: %s", err)
	}

	if rsp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(rsp.Body)
		return []int{}, fmt.Errorf("Expected 200, got %d: %s", rsp.StatusCode, b)
	}

	stories := []int{}
	if err := json.NewDecoder(rsp.Body).Decode(&stories); err != nil {
		return []int{}, fmt.Errorf("Client.TopStories: %s", err)
	}

	return stories, nil
}

func (c *baseClient) Story(id int) (*Story, error) {
	rsp, err := c.h.Get(fmt.Sprintf(itemAPI, id))

	if rsp != nil {
		defer rsp.Body.Close()
	}

	if err != nil {
		return nil, fmt.Errorf("Client.Story: %s", err)
	}

	if rsp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(rsp.Body)
		return nil, fmt.Errorf("Expected 200, got %d: %s", rsp.StatusCode, b)
	}

	var story Story
	if err := json.NewDecoder(rsp.Body).Decode(&story); err != nil {
		return nil, fmt.Errorf("Client.Story: %s", err)
	}

	return &story, nil
}

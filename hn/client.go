package hn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
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
	return []int{14779881, 14780709, 14778977, 14778685, 14780159, 14779509, 14778335}, nil
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

var cache = map[int]*Story{
	14779881: &Story{
		ID:    14779881,
		Score: 571,
		Title: "Apache Foundation bans projects from using React's “BSD+Patent” Code",
	},
	14780709: &Story{
		ID:    14780709,
		Score: 35,
		Title: "The Deal on the Table (1994)",
	},
	14778977: &Story{
		ID:    14778977,
		Score: 177,
		Title: "A 32-year-old state senator is trying to get patent trolls out of Mass (techcrunch.com)",
	},
	14778685: &Story{
		ID:    14778685,
		Score: 279,
		Title: "Monolith First (2015)",
	},
	14780159: &Story{
		ID:    14780159,
		Score: 71,
		Title: "A deep dive into Multicore OCaml garbage collector",
	},
	14779509: &Story{
		ID:    14779509,
		Score: 101,
		Title: "Why MAC address randomization is not enough [pdf]",
	},
	14778335: &Story{
		ID:    14778335,
		Score: 207,
		Title: "Tokyo street fashion and culture: 1980 – 2017 ",
	},
}

var cm sync.Mutex

func (c *baseClient) Story(id int) (*Story, error) {
	cm.Lock()
	defer cm.Unlock()

	return cache[id], nil
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

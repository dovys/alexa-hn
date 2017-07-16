package alexahn

import (
	"context"
	"fmt"
	"sync"

	"sort"

	"github.com/dovys/alexa-hn/hn"
)

const (
	frontPageLimit = 30
	storiesToRead  = 5
)

type AlexaResponse struct {
	Speech string
	Card   *struct {
		Title   string
		Content string
	}
}

type Service interface {
	ReadTopStories(context.Context) (*AlexaResponse, error)
}

func NewService(hn hn.Client) Service {
	return &service{hn: hn}
}

type service struct {
	hn hn.Client
}

func (s *service) ReadTopStories(context.Context) (*AlexaResponse, error) {
	top, err := s.hn.TopStories()
	if err != nil {
		return nil, err
	}

	if len(top) > frontPageLimit {
		top = top[0:5] // todo: use frontPageLimit
	}

	wg := &sync.WaitGroup{}
	stories := make([]*hn.Story, 0, len(top))
	rcv := make(chan *hn.Story, len(top))

	for _, id := range top {
		wg.Add(1)
		go func(w chan<- *hn.Story, wg *sync.WaitGroup, id int) {
			defer wg.Done()
			story, err := s.hn.Story(id)

			// todo: err chan
			if err != nil {
				return
			}
			w <- story
		}(rcv, wg, id)
	}

	wg.Wait()
	close(rcv)

	for st := range rcv {
		stories = append(stories, st)
	}

	sort.Slice(stories, func(i, j int) bool {
		return stories[i].Score > stories[j].Score
	})

	if len(stories) > storiesToRead {
		stories = stories[0:storiesToRead]
	}

	return buildAlexaResponse(stories), nil
}

func buildAlexaResponse(stories []*hn.Story) *AlexaResponse {
	if len(stories) == 0 {
		return &AlexaResponse{Speech: "Could not read the stories"}
	}

	text := fmt.Sprintf("Reading top %d stories. ", len(stories))
	for i := 0; i < len(stories); i++ {
		text += fmt.Sprintf("%d points. %s. . . . ", stories[i].Score, stories[i].Title)
	}

	return &AlexaResponse{Speech: text}
}

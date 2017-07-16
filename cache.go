package alexahn

import (
	"sync"
	"time"
)

// todo: use app engine memcached
func Cache(svc Service, ttl time.Duration) Service {
	return &cachingDecorator{svc: svc, ttl: ttl, mutex: &sync.Mutex{}}
}

type cachingDecorator struct {
	svc Service

	ttl      time.Duration
	updated  time.Time
	response *AlexaResponse
	mutex    *sync.Mutex
}

func (s *cachingDecorator) ReadTopStories() (*AlexaResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.response != nil && s.updated.Add(s.ttl).After(time.Now()) {
		return s.response, nil
	}

	c, err := s.svc.ReadTopStories()

	if err != nil {
		return nil, err
	}

	s.response = c
	s.updated = time.Now()

	return c, nil
}

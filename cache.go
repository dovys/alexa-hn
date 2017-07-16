package alexahn

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/appengine/memcache"
)

// todo: use app engine memcached
func InMemoryCache(svc Service, ttl time.Duration) Service {
	return &staticCacheDecorator{svc: svc, ttl: ttl, mutex: &sync.Mutex{}}
}

func MemcachedCache(svc Service, ttl time.Duration) Service {
	return &memcachedDecorator{svc: svc, ttl: ttl}
}

type memcachedDecorator struct {
	svc Service
	ttl time.Duration
}

func (s *memcachedDecorator) ReadTopStories(ctx context.Context) (*AlexaResponse, error) {
	i, err := memcache.Get(ctx, "top_stories")

	fmt.Println("memcache.Get", err)

	if err != nil && err != memcache.ErrCacheMiss {
		fmt.Println(err)
		return s.svc.ReadTopStories(ctx)
	}

	if err == memcache.ErrCacheMiss || i == nil || i.Object == nil {
		r, err := s.svc.ReadTopStories(ctx)

		if err != nil {
			return nil, err
		}

		fmt.Println("memcache.Add", memcache.Add(ctx, &memcache.Item{
			Key:        "top_stories",
			Object:     *r,
			Expiration: s.ttl,
		}))

		return r, nil
	}

	fmt.Printf("%+v\n%s\n%s\n", i.Object, i.Value, i.Key)

	r := i.Object.(AlexaResponse)

	return &r, nil
}

type staticCacheDecorator struct {
	svc Service

	ttl      time.Duration
	updated  time.Time
	response *AlexaResponse
	mutex    *sync.Mutex
}

func (s *staticCacheDecorator) ReadTopStories(ctx context.Context) (*AlexaResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.response != nil && s.updated.Add(s.ttl).After(time.Now()) {
		return s.response, nil
	}

	c, err := s.svc.ReadTopStories(ctx)

	if err != nil {
		return nil, err
	}

	s.response = c
	s.updated = time.Now()

	return c, nil
}

package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type slideWindowLimitter struct {
	ctx          context.Context
	timeDuration int
	rate         int
	tokens       sync.Map
}

type slideWindowClient struct {
	expireTime time.Time
	tokens     chan struct{}
}

func NewSlideWindowLimitter(timeduration int, rate int) *slideWindowLimitter {
	return &slideWindowLimitter{
		timeDuration: timeduration,
		rate:         rate,
		tokens:       sync.Map{},
	}
}

func (l *slideWindowLimitter) HandleRequest(ip string) {
	v, ok := l.GetToken(ip)
	if !ok {
		l.initialIp(ip)
		l.HandleRequest(ip)
		return
	}

	if v.expireTime.After(time.Now()) {
		select {
		case <-v.tokens:
			fmt.Println("Request Accepted!")
		default:
			fmt.Println("Request Denied!")
		}
		return
	}

	v.expireTime = time.Now().Add(time.Duration(l.timeDuration) * time.Second)
	l.refileTokens(ip)

	l.HandleRequest(ip)
}

func (l *slideWindowLimitter) initialIp(ip string) {
	l.tokens.Store(ip, &slideWindowClient{
		expireTime: time.Now().Add(time.Duration(l.timeDuration) * time.Second),
		tokens:     make(chan struct{}, l.rate),
	})

	l.refileTokens(ip)
}

func (l *slideWindowLimitter) refileTokens(ip string) {
	tokensList, ok := l.GetToken(ip)
	if !ok {
		return
	}
	for i := 0; i < l.rate; i++ {
		select {
		case tokensList.tokens <- struct{}{}:
		default:
		}
	}
}

func (s *slideWindowLimitter) GetToken(key string) (*slideWindowClient, bool) {
	value, ok := s.tokens.Load(key)
	if !ok {
		return nil, false
	}
	client, valid := value.(*slideWindowClient)
	return client, valid
}

func (s *slideWindowLimitter) SetToken(key string, client *slideWindowClient) {
	s.tokens.Store(key, client)
}

func (s *slideWindowLimitter) DeleteToken(key string) {
	s.tokens.Delete(key)
}

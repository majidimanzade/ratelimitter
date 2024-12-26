package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

type Ratelimit struct {
	TokensByIp   sync.Map
	RefileRate   time.Duration
	RefileCount  int
	Stop         chan struct{}
	Capacity     int
	TruncateTime time.Duration
}

func NewTokenRateLimitter(capacity int) *Ratelimit {
	limitter := &Ratelimit{
		// TokensByIp:   make(map[string]chan struct{}, capacity),
		TokensByIp:   sync.Map{},
		RefileRate:   time.Duration(time.Second),
		RefileCount:  2,
		Stop:         make(chan struct{}),
		Capacity:     capacity,
		TruncateTime: time.Duration(time.Second * 3),
	}

	go limitter.refileTokens()
	go limitter.truncateIps()

	return limitter
}
func (l *Ratelimit) refileTokens() {
	ticker := time.NewTicker(l.RefileRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.TokensByIp.Range(func(k, v interface{}) bool {
				ch, ok := v.(chan struct{})
				if !ok {
					return true
				}

				for i := 0; i < l.RefileCount; i++ {
					select {
					case ch <- struct{}{}:
					default:
					}
				}
				return true
			})
		case <-l.Stop:
			return
		}
	}
}

func (l *Ratelimit) truncateIps() {
	ticker := time.NewTicker(l.TruncateTime)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.TokensByIp.Range(func(k, v interface{}) bool {
				ch, ok := v.(chan struct{})
				if !ok {
					return false
				}

				select {
				case ch <- struct{}{}:
				default:
					l.TokensByIp.Delete(k)
				}
				return true
			})
		}
	}
}

func (l *Ratelimit) RefileTokensNewIp(ip string) {
	c := make(chan struct{}, l.Capacity)
	for i := 0; i < l.Capacity; i++ {
		c <- struct{}{}
	}
	l.TokensByIp.Store(ip, c)
}

func (l *Ratelimit) HandleRequest(v int, ip string) {
	_, ok := l.TokensByIp.Load(ip)
	if !ok {
		l.RefileTokensNewIp(ip)
	}

	ch, _ := l.TokensByIp.Load(ip)
	chh := ch.(chan struct{})
	select {
	case <-chh:
		fmt.Printf(" \n Request Accepted ip:%s %d\n ", ip, v)
	default:
		fmt.Printf("\n Request Denied ip:%s %d\n ", ip, v)
	}
}

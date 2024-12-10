package ratelimit

import "time"

type Ratelimit struct {
	TokensByIp  map[string]chan struct{}
	RefileRate  time.Duration
	RefileCount int
	Stop        chan struct{}
	Capacity    int
}

func NewTokenRateLimitter(capacity int) Ratelimit {
	limitter := Ratelimit{
		TokensByIp:  make(map[string]chan struct{}, capacity),
		RefileRate:  time.Duration(time.Second),
		RefileCount: 2,
		Stop:        make(chan struct{}),
		Capacity:    capacity,
	}

	go limitter.refileTokens()

	return limitter
}

func (l Ratelimit) refileTokens() {
	ticker := time.NewTicker(l.RefileRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, ch := range l.TokensByIp {
				for i := 0; i < l.RefileCount; i++ {
					select {
					case ch <- struct{}{}:
					default:
					}
				}
			}
		case <-l.Stop:
			return
		}

	}
}

func (l Ratelimit) RefileTokensNewIp(ip string) {
	c := make(chan struct{}, l.Capacity)
	for i := 0; i < l.Capacity; i++ {
		c <- struct{}{}
	}
	l.TokensByIp[ip] = c
}
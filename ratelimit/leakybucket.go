package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

type LeakyBucketLimitter struct {
	QueueSize int8
	MapReset  time.Duration
	Queue     sync.Map
}

func NewLeakyBucketLimitter(size int8, mapReset time.Duration) *LeakyBucketLimitter {
	limitter := LeakyBucketLimitter{
		QueueSize: size,
		MapReset:  mapReset,
		Queue:     sync.Map{},
	}

	return &limitter
}

func (l *LeakyBucketLimitter) HandleRequest(id int, ip string) {
	ch, ok := l.isKeyExist(ip)
	if !ok {
		l.initIpQueue(ip)
	}

	select {
	case <-ch:
		fmt.Printf("\nRequest Accepted! %s %d\n", ip, id)
		time.Sleep(time.Second)
		ch <- struct{}{}
	default:
		fmt.Printf("\nRequest Denied! %s %d\n", ip, id)
	}
}

func (l *LeakyBucketLimitter) initIpQueue(ip string) {
	l.Queue.Store(ip, make(chan struct{}, l.QueueSize))
	ch, ok := l.isKeyExist(ip)
	if !ok {
		return
	}

	for i := 0; i < int(l.QueueSize); i++ {
		ch <- struct{}{}
	}
}

func (l *LeakyBucketLimitter) isKeyExist(key string) (chan struct{}, bool) {
	v, ok := l.Queue.Load(key)
	if !ok {
		return nil, false
	}

	ch, ok := v.(chan struct{})
	if !ok {
		fmt.Println("Error: Invalid channel type")
		return nil, false
	}
	return ch, true
}

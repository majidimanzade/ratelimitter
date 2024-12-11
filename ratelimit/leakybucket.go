package ratelimit

import (
	"fmt"
	"time"
)

type LeakyBucketLimitter struct {
	QueueSize int8
	MapReset  time.Duration
	Queue     map[string]chan struct{}
}

func NewLeakyBucketLimitter(size int8, mapReset time.Duration) LeakyBucketLimitter {
	limitter := LeakyBucketLimitter{
		QueueSize: size,
		MapReset:  mapReset,
		Queue:     make(map[string]chan struct{}),
	}

	return limitter
}

func (l LeakyBucketLimitter) HandleRequest(id int, ip string) {
	if _, ok := l.Queue[ip]; !ok {
		l.initIpQueue(ip)
	}

	select {
	case <-l.Queue[ip]:
		fmt.Printf("\nRequest Accepted! %s %d\n", ip, id)
		time.Sleep(time.Second)
		l.Queue[ip] <- struct{}{}
	default:
		fmt.Printf("\nRequest Denied! %s %d\n", ip, id)
	}
}

func (l LeakyBucketLimitter) initIpQueue(ip string) {
	l.Queue[ip] = make(chan struct{}, l.QueueSize)
	fmt.Println(l.QueueSize)
	for i := 0; i < int(l.QueueSize); i++ {
		l.Queue[ip] <- struct{}{}
	}
}

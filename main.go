package main

import (
	"context"
	"fmt"
	"myproject/ratelimit"
	"time"
)

func main() {
	// leackyBucketLimitter()
	tokenBucketLimitter()
	// slideWindowLimitter()
}

func leackyBucketLimitter() {
	limitter := ratelimit.NewLeakyBucketLimitter(2, time.Second*5)

	for i := 0; i < 11; i++ {
		if i == 2 || i == 3 {
			go limitter.HandleRequest(i, "192.168.1.2")
		} else {
			go limitter.HandleRequest(i, "192.168.1.1")
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func tokenBucketLimitter() {
	limitter := ratelimit.NewTokenRateLimitter(2)

	for i := 0; i < 11; i++ {
		if i == 2 || i == 3 {
			limitter.HandleRequest(i, "192.168.1.2")
		} else {
			limitter.HandleRequest(i, "192.168.1.1")
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func slideWindowLimitter() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))

	limitter := ratelimit.NewSlideWindowLimitter(2, 5)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Session Expired!")
			return
		default:
			limitter.HandleRequest("192.168.1.1")
			time.Sleep(200 * time.Millisecond)
		}
	}
}

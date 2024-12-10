package main

import (
	"fmt"
	"myproject/ratelimit"
	"time"
)

func main() {
	limiter := ratelimit.NewTokenRateLimitter(2)

	for i := 0; i < 11; i++ {
		if i == 2 || i == 3 {
			handleRequest(i, limiter, "192.168.1.2")
		} else {
			handleRequest(i, limiter, "192.168.1.1")
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func handleRequest(v int, limitter ratelimit.Ratelimit, ip string) {
	if _, ok := limitter.TokensByIp[ip]; !ok {
		limitter.RefileTokensNewIp(ip)
	}

	select {
	case <-limitter.TokensByIp[ip]:
		fmt.Printf(" \n Request Accepted ip:%s %d\n ", ip, v)
	default:
		fmt.Printf("\n Request Denied ip:%s %d\n ", ip, v)
	}
}

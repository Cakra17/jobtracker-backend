package ratelimiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int64 	// max num of token
	tokens     int64	// current number of token
	refillRate int64	// token refilled per second
	lastRefill time.Time // last time bucket refilled
	mu sync.Mutex // make sure thread safe
}

func NewTokenBucket(capacity, refillRate int64) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		tokens: capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) AllowN(n int64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.Refill()

	if tb.tokens  >= n {
		tb.tokens -= n
		return true
	}
	return false
}

func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1)
}

func (tb *TokenBucket) Refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	tokenToAdd := int64(elapsed.Seconds()) * tb.refillRate

	if tokenToAdd > 0 {
		tb.tokens = min(tb.capacity, tokenToAdd + tb.tokens)
		tb.lastRefill = now
	}
}

func (tb *TokenBucket) AvailableToken() int64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.Refill()
	return tb.tokens
}
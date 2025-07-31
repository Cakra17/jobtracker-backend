package ratelimiter

import (
	"errors"
	"sync"
	"time"
)

type RateLimiter struct {
	buckets    map[string]*TokenBucket
	capacity   int64
	refillRate int64
	mu         sync.RWMutex
}

func NewRateLimiter(capacity, refillRate int64) RateLimiter {
	return RateLimiter{
		buckets: make(map[string]*TokenBucket),
		capacity: capacity,
		refillRate: refillRate,
	}
}

func (rl *RateLimiter) getBucket(key string) *TokenBucket {
	rl.mu.RLock()
	bucket, ok := rl.buckets[key]
	rl.mu.RUnlock()

	if ok {
		return bucket
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if bucket, ok := rl.buckets[key]; ok {
		return bucket
	}

	bucket = NewTokenBucket(rl.capacity, rl.refillRate)
	return bucket
}

func (rl *RateLimiter) AllowN(key string,n int64) bool {
	bucket := rl.getBucket(key)
	return bucket.AllowN(n)
}

func (rl *RateLimiter) Allow(key string) bool {
	bucket := rl.getBucket(key)
	return bucket.Allow()
}

func (rl *RateLimiter) Cleanup(maxAge time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	for key, bucket := range rl.buckets {
		bucket.mu.Lock()
		if bucket.lastRefill.Before(cutoff) {
			delete(rl.buckets, key)
		}
		bucket.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(ip string) error {
	if !rl.Allow(ip) {
		return errors.New("request got suspend")
	}
	return nil
}
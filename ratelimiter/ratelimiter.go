package ratelimiter

import (
	"sync"
	"time"
)

var globalLimiter = newRateLimiter(10.0, 20.0)

// RateLimiter implements token bucket algorithm for rate limiting
type RateLimiter struct {
	tokens     map[string]float64
	lastUpdate map[string]time.Time
	rate       float64
	capacity   float64
	mu         sync.Mutex
	ipTimeout  time.Duration
	cleanupInt time.Duration
	quit       chan struct{}
}

// newRateLimiter creates a new rate limiter with specified rate and capacity
// Optional custom timeout and cleanup interval (set to 0 to use default)
func newRateLimiter(rate float64, capacity float64, options ...time.Duration) *RateLimiter {
	ipTimeout := 30 * time.Minute
	cleanupInt := 5 * time.Minute

	if len(options) > 0 && options[0] > 0 {
		ipTimeout = options[0]
	}
	if len(options) > 1 && options[1] > 0 {
		cleanupInt = options[1]
	}

	rl := &RateLimiter{
		tokens:     make(map[string]float64),
		lastUpdate: make(map[string]time.Time),
		rate:       rate,
		capacity:   capacity,
		ipTimeout:  ipTimeout,
		cleanupInt: cleanupInt,
		quit:       make(chan struct{}),
	}

	go rl.cleanup()

	return rl
}

// Stop stops the cleanup routine
func (rl *RateLimiter) Stop() {
	close(rl.quit)
}

// cleanup periodically removes expired IPs from memory
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInt)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for ip, lastSeen := range rl.lastUpdate {
				if now.Sub(lastSeen) > rl.ipTimeout {
					delete(rl.tokens, ip)
					delete(rl.lastUpdate, ip)
				}
			}
			rl.mu.Unlock()
		case <-rl.quit:
			return
		}
	}
}

// Allow checks if a request from the given IP should be allowed
// Returns true if allowed, false if rate limited
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	if _, exists := rl.tokens[ip]; !exists {
		rl.tokens[ip] = rl.capacity
		rl.lastUpdate[ip] = now
		return true
	}

	// Calculate tokens to add based on time elapsed
	elapsed := now.Sub(rl.lastUpdate[ip]).Seconds()
	rl.tokens[ip] = min(rl.capacity, rl.tokens[ip]+(elapsed*rl.rate))
	rl.lastUpdate[ip] = now

	if rl.tokens[ip] >= 1.0 {
		rl.tokens[ip] -= 1.0
		return true
	}

	return false
}

// GetStats returns the current token count and last update time for an IP
func (rl *RateLimiter) GetStats(ip string) (float64, bool, time.Time) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	tokens, exist := rl.tokens[ip]
	return tokens, exist, rl.lastUpdate[ip]
}

// GetGlobalLimiter returns the global rate limiter
func GetGlobalLimiter() *RateLimiter {
	return globalLimiter
}

// ShutdownLimiter should be called when the application is shutting down
func ShutdownLimiter() {
	if globalLimiter != nil {
		globalLimiter.Stop()
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

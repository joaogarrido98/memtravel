package ratelimiter

import (
	"testing"
	"time"
)

func TestAllow_NewIP_Allows(t *testing.T) {
	rl := newRateLimiter(1.0, 2.0)
	defer rl.Stop()

	ip := "192.168.1.1"
	allowed := rl.Allow(ip)

	if !allowed {
		t.Errorf("Expected new IP to be allowed")
	}
}

func TestAllow_RateLimitExceeded(t *testing.T) {
	rl := newRateLimiter(0.0, 2.0)
	defer rl.Stop()

	ip := "192.168.1.2"

	if !rl.Allow(ip) {
		t.Fatal("1st request should be allowed")
	}
	if !rl.Allow(ip) {
		t.Fatal("2nd request should be allowed")
	}
	if !rl.Allow(ip) {
		t.Error("3rd request should be allowed")
	}
	if rl.Allow(ip) {
		t.Error("4th request should be rate-limited (no tokens left)")
	}
}

func TestAllow_TokenRefill(t *testing.T) {
	rl := newRateLimiter(1.0, 2.0)
	defer rl.Stop()

	ip := "192.168.1.3"

	rl.Allow(ip) // uses 1 token
	time.Sleep(1100 * time.Millisecond)

	if !rl.Allow(ip) {
		t.Errorf("Expected token to refill after 1 second")
	}
}

func TestGetStats_ReturnsCorrectValues(t *testing.T) {
	rl := newRateLimiter(1.0, 2.0)
	defer rl.Stop()

	ip := "192.168.1.4"
	rl.Allow(ip)

	tokens, exist, lastSeen := rl.GetStats(ip)
	if !exist {
		t.Errorf("Expected IP to exist in stats")
	}
	if tokens > rl.capacity || tokens < 0 {
		t.Errorf("Token count out of bounds: %f", tokens)
	}
	if lastSeen.IsZero() {
		t.Errorf("Expected non-zero lastSeen time")
	}
}

func TestCleanup_RemovesStaleEntries(t *testing.T) {
	rl := newRateLimiter(1.0, 2.0, 1*time.Second, 500*time.Millisecond)
	defer rl.Stop()

	ip := "192.168.1.5"
	rl.Allow(ip)

	time.Sleep(2 * time.Second) // Wait for cleanup to remove stale entry

	_, exist, _ := rl.GetStats(ip)
	if exist {
		t.Errorf("Expected IP to be removed by cleanup")
	}
}

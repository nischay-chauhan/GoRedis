package middleware

import (
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	visitors     *sync.Map
	r            int // max requests per second
	b            int // max burst size
	maxVisitors  int32
	currentCount int32
	stopCleanup  chan struct{}
	stopOnce     sync.Once
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen atomic.Int64
}

func NewRateLimiter(r, b int, cleanupInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors:     &sync.Map{},
		r:            r,
		b:            b,
		maxVisitors:  0,
		stopCleanup:  make(chan struct{}),
		currentCount: 0,
	}

	go rl.cleanupVisitors(cleanupInterval)
	return rl
}

func (rl *RateLimiter) Stop() {
	rl.stopOnce.Do(func() {
		close(rl.stopCleanup)
	})
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		if ip == "" {
			http.Error(w, "Unable to identify IP", http.StatusBadRequest)
			return
		}

		log.Printf("[rate] incoming %s %s ip=%s", r.Method, r.URL.Path, ip)

		limiter, remaining, resetTime := rl.getVisitor(ip)
		if limiter == nil {
			log.Printf("[rate] deny (too many visitors) ip=%s", ip)
			http.Error(w, "Too many users, please try again later", http.StatusServiceUnavailable)
			return
		}

		if !limiter.Allow() {
			retryAfter := int(time.Until(resetTime).Seconds())
			if retryAfter < 0 {
				retryAfter = 1
			}

			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.r))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

			// Log deny decision
			log.Printf("[rate] deny 429 %s %s ip=%s rem=0 reset=%d retry=%d", r.Method, r.URL.Path, ip, resetTime.Unix(), retryAfter)
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Set rate limit headers
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.r))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) getVisitor(ip string) (*rate.Limiter, int, time.Time) {
    now := time.Now()
    resetTime := now.Truncate(time.Second).Add(time.Second)
    if val, ok := rl.visitors.Load(ip); ok {
        v := val.(*visitor)
        v.lastSeen.Store(now.UnixNano())
        remaining := int(v.limiter.Tokens())
        return v.limiter, remaining, resetTime
    }

    // Create new visitor
    limiter := rate.NewLimiter(rate.Limit(rl.r), rl.b)
    v := &visitor{
        limiter:  limiter,
		lastSeen: atomic.Int64{},
	}
	v.lastSeen.Store(now.UnixNano())

	if existing, loaded := rl.visitors.LoadOrStore(ip, v); loaded {
		exV := existing.(*visitor)
		exV.lastSeen.Store(now.UnixNano())
		remaining := int(exV.limiter.Tokens())
		return exV.limiter, remaining, resetTime
	}

	if rl.maxVisitors > 0 {
		if atomic.AddInt32(&rl.currentCount, 1) > rl.maxVisitors {
			rl.visitors.Delete(ip)
			atomic.AddInt32(&rl.currentCount, -1)
			return nil, 0, time.Time{}
		}
	}

	return limiter, rl.b, resetTime
}

func (rl *RateLimiter) cleanupVisitors(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.visitors.Range(func(key, value interface{}) bool {
				v := value.(*visitor)
				last := time.Unix(0, v.lastSeen.Load())
				if time.Since(last) > 3*time.Minute {
					rl.visitors.Delete(key)
					if rl.maxVisitors > 0 {
						atomic.AddInt32(&rl.currentCount, -1)
					}
				}
				return true
			})
		case <-rl.stopCleanup:
			return
		}
	}
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}

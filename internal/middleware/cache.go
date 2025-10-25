package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheConfig struct {
	DefaultTTL time.Duration
	SkipCacheHeader string
	CacheControl bool
}

type cacheEntry struct {
	Status  int
	Headers http.Header
	Body    []byte
}

func NewCache(rdb *redis.Client, config CacheConfig) func(http.Handler) http.Handler {
	if config.DefaultTTL == 0 {
		config.DefaultTTL = 5 * time.Minute
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			if config.SkipCacheHeader != "" && r.Header.Get(config.SkipCacheHeader) != "" {
				next.ServeHTTP(w, r)
				return
			}

			key := generateCacheKey(r)

			cached, err := getFromCache(rdb, key)
			if err == nil && cached != nil {
				// Serve from cache
				for k, v := range cached.Headers {
					w.Header()[k] = v
				}
				w.WriteHeader(cached.Status)
				w.Write(cached.Body)
				return
			}

			rw := &responseRecorder{ResponseWriter: w, body: []byte{}}

			next.ServeHTTP(rw, r)

			if rw.status >= 200 && rw.status < 300 {
				ttl := config.DefaultTTL
				if config.CacheControl {
					if cacheControl := w.Header().Get("Cache-Control"); cacheControl != "" {
						if maxAge, ok := parseMaxAge(cacheControl); ok {
							ttl = time.Duration(maxAge) * time.Second
						}
					}
				}

				entry := &cacheEntry{
					Status:  rw.status,
					Headers: w.Header().Clone(),
					Body:    rw.body,
				}

				if w.Header().Get("Cache-Control") == "" {
					w.Header().Set("Cache-Control", 
						"public, max-age="+strconv.Itoa(int(ttl.Seconds())))
				}

				go setInCache(rdb, key, entry, ttl)
			}

			w.WriteHeader(rw.status)
			w.Write(rw.body)
		})
	}
}

func generateCacheKey(r *http.Request) string {
	hash := sha256.New()
	hash.Write([]byte(r.Method))
	hash.Write([]byte(r.URL.Path))
	hash.Write([]byte(r.URL.RawQuery))
	hash.Write([]byte(r.Header.Get("Accept")))
	hash.Write([]byte(r.Header.Get("Accept-Language")))
	return "cache:" + hex.EncodeToString(hash.Sum(nil))
}

func parseMaxAge(cacheControl string) (int, bool) {
	parts := strings.Split(cacheControl, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "max-age=") {
			if age, err := strconv.Atoi(part[8:]); err == nil {
				return age, true
			}
		}
	}
	return 0, false
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return len(b), nil
}

func getFromCache(rdb *redis.Client, key string) (*cacheEntry, error) {
	ctx := context.Background()
	data, err := rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var entry cacheEntry
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

func setInCache(rdb *redis.Client, key string, entry *cacheEntry, ttl time.Duration) error {
	ctx := context.Background()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(entry); err != nil {
		return err
	}

	return rdb.Set(ctx, key, buf.Bytes(), ttl).Err()
}

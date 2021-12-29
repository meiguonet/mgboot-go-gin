package taskx

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/meiguonet/mgboot-go-dal/poolx"
	"github.com/meiguonet/mgboot-go-fiber/cachex"
	"sync"
	"time"
)

type redismqNormalTaskHandler struct {
}

func (h *redismqNormalTaskHandler) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := poolx.GetRedisConnection(ctx)

	if err != nil {
		return
	}

	defer conn.Close()
	cacheKey := cachex.CacheKeyRedismqNormal()
	payloads := make([]string, 0)
	wg := sync.WaitGroup{}
	wg.Add(10)
	mu := sync.Mutex{}

	for i := 1; i <= 10; i++ {
		go func() {
			defer wg.Done()
			payload, _ := redis.String(conn.Do("LPOP", cacheKey))

			if payload == "" {
				return
			}

			mu.Lock()
			payloads = append(payloads, payload)
			mu.Unlock()
		}()
	}

	wg.Wait()

	if len(payloads) < 1 {
		return
	}

	wg = sync.WaitGroup{}
	wg.Add(len(payloads))

	for _, payload := range payloads {
		go func(payload string) {
			defer wg.Done()
			RunMqTask(payload)
		}(payload)
	}

	wg.Wait()
}

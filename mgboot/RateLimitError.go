package mgboot

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type RateLimitError struct {
	total      int
	remaining  int
	retryAfter string
}

func NewRateLimitError(data map[string]interface{}) RateLimitError {
	var total int

	if n1, ok := data["total"].(int); ok && n1 > 0 {
		total = n1
	}

	var remaining int

	if n1, ok := data["remaining"].(int); ok && n1 > 0 {
		remaining = n1
	}

	var retryAfter string

	if s1, ok := data["retryAfter"].(string); ok && s1 != "" {
		retryAfter = s1
	}

	return RateLimitError{
		total:      total,
		remaining:  remaining,
		retryAfter: retryAfter,
	}
}

func (ex RateLimitError) Error() string {
	return "rate limit exceed"
}

func (ex RateLimitError) Total() int {
	return ex.total
}

func (ex RateLimitError) Remaining() int {
	return ex.remaining
}

func (ex RateLimitError) RetryAfter() string {
	return ex.retryAfter
}

func (ex RateLimitError) AddSpecifyHeaders(ctx *gin.Context) {
	ctx.Header("X-Ratelimit-Limit", fmt.Sprintf("%d", ex.Total()))
	ctx.Header("X-Ratelimit-Remaining", fmt.Sprintf("%d", ex.Remaining()))

	if ex.RetryAfter() != "" {
		ctx.Header("Retry-After", ex.RetryAfter())
	}
}

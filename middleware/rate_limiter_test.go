package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type clockMock struct {
	nowMock    func() int64
	tIndex     int
	timestamps []int64
}

func createClock(timestampsSeconds ...int64) clockMock {
	return clockMock{
		timestamps: timestampsSeconds,
		tIndex:     0,
	}
}

func (t *clockMock) nowNano() int64 {
	if t.tIndex > len(t.timestamps) {
		panic("error")
	}
	timestamp := t.timestamps[t.tIndex]
	t.tIndex++
	return timestamp * time.Second.Nanoseconds()
}

func TestRateLimiter_Accept(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		clMock := createClock(1, 2, 3)
		rateLimiter := newRateLimiterWithClock(3, 3, &clMock)
		httpRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpRecorder)
		ctx.Request, _ = http.NewRequestWithContext(ctx, "GET", "/message/", nil)
		ctx.Request.Header.Add("userId", "pepe")

		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
	})
	t.Run("success", func(t *testing.T) {
		clMock := createClock(1, 2)
		rateLimiter := newRateLimiterWithClock(2, 1, &clMock)
		httpRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpRecorder)
		ctx.Request, _ = http.NewRequestWithContext(ctx, "GET", "/message/", nil)
		ctx.Request.Header.Add("userId", "pepe")

		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
	})
	t.Run("fail", func(t *testing.T) {
		clMock := createClock(1, 2, 3)
		rateLimiter := newRateLimiterWithClock(2, 3, &clMock)
		httpRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpRecorder)
		ctx.Request, _ = http.NewRequestWithContext(ctx, "GET", "/message/", nil)
		ctx.Request.Header.Add("userId", "pepe")

		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusTooManyRequests)
	})
	t.Run("fail2", func(t *testing.T) {
		clMock := createClock(1, 1, 2, 2, 3, 3)
		rateLimiter := newRateLimiterWithClock(4, 5, &clMock)
		httpRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpRecorder)
		ctx.Request, _ = http.NewRequestWithContext(ctx, "GET", "/message/", nil)
		ctx.Request.Header.Add("userId", "pepe")
		ctx2, _ := gin.CreateTestContext(httpRecorder)
		ctx2.Request, _ = http.NewRequestWithContext(ctx, "GET", "/message/", nil)
		ctx2.Request.Header.Add("userId", "pepe2")

		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx2)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx2)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx2)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
	})
	t.Run("fail2", func(t *testing.T) {
		clMock := createClock(1, 1, 2, 2, 3, 3, 4)
		rateLimiter := newRateLimiterWithClock(3, 5, &clMock)
		httpRecorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(httpRecorder)
		ctx.Request, _ = http.NewRequestWithContext(ctx, "GET", "/message/", nil)
		ctx.Request.Header.Add("userId", "pepe")
		ctx2, _ := gin.CreateTestContext(httpRecorder)
		ctx2.Request, _ = http.NewRequestWithContext(ctx, "GET", "/message/", nil)
		ctx2.Request.Header.Add("userId", "pepe2")

		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx2)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx2)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx2)
		assert.Equal(t, httpRecorder.Code, http.StatusOK)
		rateLimiter.Accept(ctx2)
		assert.Equal(t, httpRecorder.Code, http.StatusTooManyRequests)
	})
}

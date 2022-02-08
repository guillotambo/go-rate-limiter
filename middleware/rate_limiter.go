package middleware

import (
	"container/list"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Limiter interface {
	Accept() gin.HandlerFunc
}

//La idea de esta interface, es para test
type clock interface {
	nowNano() int64
}

type SystemTimeClock struct{}

func (t *SystemTimeClock) nowNano() int64 {
	return time.Now().UnixNano()
}

type (
	RateLimiter struct {
		maxRequestAllowedPerUser int
		timeFrameInNano          int64
		timestampsByUser         map[string]*list.List //La lista va a estar siempre ordeanado, porque usa el time.Now del clock
		clock                    clock
	}
)

func NewRateLimiter(maxRequestAllowedPerUser int, timeFrameInSeconds int64) RateLimiter {
	return newRateLimiterWithClock(
		maxRequestAllowedPerUser,
		timeFrameInSeconds,
		&SystemTimeClock{})
}

func newRateLimiterWithClock(maxRequestAllowedPerUser int, timeFrameInSeconds int64, clock clock) RateLimiter {
	return RateLimiter{
		maxRequestAllowedPerUser: maxRequestAllowedPerUser,
		timeFrameInNano:          timeFrameInSeconds * time.Second.Nanoseconds(),
		timestampsByUser:         map[string]*list.List{},
		clock:                    clock,
	}
}

func (rateLimiter *RateLimiter) Accept(ctx *gin.Context) {
	userId := ctx.GetHeader("userId")
	if userId == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "user Id is not valid")
		return
	}
	currentTimestamp := rateLimiter.clock.nowNano()

	if rateLimiter.timestampsByUser[userId] == nil {
		rateLimiter.timestampsByUser[userId] = &list.List{}
	}

	userTimestamps := rateLimiter.timestampsByUser[userId]

	rateLimiter.cleanOldRequests(userTimestamps, currentTimestamp)

	if !(userTimestamps.Len() < rateLimiter.maxRequestAllowedPerUser) {
		ctx.AbortWithStatusJSON(http.StatusTooManyRequests, "Try again later ;(")
		return
	}

	userTimestamps.PushBack(currentTimestamp)

	ctx.Next()
}

func (rateLimiter *RateLimiter) cleanOldRequests(userTimestamps *list.List, now int64) {
	element := userTimestamps.Front()
	oldestPossibleTime := now - rateLimiter.timeFrameInNano
	for i := 0; i < userTimestamps.Len(); i++ {
		value := element.Value.(int64)
		next := element.Next()
		if value <= oldestPossibleTime {
			userTimestamps.Remove(element)
		}
		element = next
	}
}

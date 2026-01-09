package tcp

import (
	"errors"
)

var ErrLimitReached = errors.New("max connections reached")

// ConnLimiter: 세마포어 패턴을 이용한 연결 제한기
type ConnLimiter struct {
	sem chan struct{}
}

func NewConnLimiter(max int) *ConnLimiter {
	return &ConnLimiter{
		sem: make(chan struct{}, max),
	}
}

func (l *ConnLimiter) Acquire() bool {
	select {
	case l.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

func (l *ConnLimiter) Release() {
	<-l.sem
}
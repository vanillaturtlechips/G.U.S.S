package tcp

import (
	"expvar"
)

// TCP 관련 전역 메트릭 정의
var (
	CurrentConns = expvar.NewInt("tcp_current_connections")
	TotalConns   = expvar.NewInt("tcp_total_connections")
	TimeoutErrs  = expvar.NewInt("tcp_timeout_errors")
)

// AddConn: 연결 시작 시 메트릭 업데이트
func AddConn() {
	CurrentConns.Add(1)
	TotalConns.Add(1)
}

// RemoveConn: 연결 종료 시 메트릭 업데이트
func RemoveConn() {
	CurrentConns.Add(-1)
}
package tcp

import (
	"net"
	"time"
)

// SetKeepAlive: TCP KeepAlive 설정을 적용합니다.
func SetKeepAlive(conn net.Conn, d time.Duration) {
	if tc, ok := conn.(*net.TCPConn); ok {
		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(d)
	}
}

// ApplyDeadlines: 읽기/쓰기 타임아웃을 한 번에 적용합니다.
func ApplyDeadlines(conn net.Conn, read, write time.Duration) {
	now := time.Now()
	if read > 0 {
		conn.SetReadDeadline(now.Add(read))
	}
	if write > 0 {
		conn.SetWriteDeadline(now.Add(write))
	}
}
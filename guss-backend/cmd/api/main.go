package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// MySQL 드라이버
	_ "github.com/go-sql-driver/mysql"

	"guss-backend/internal/algo"
	"guss-backend/internal/api"
	"guss-backend/internal/repository"
	"guss-backend/pkg/tcp"
)

// HTTP 응답 규격 준수를 위한 ResponseWriter
type mockResponseWriter struct {
	conn        net.Conn
	header      http.Header
	wroteHeader bool
}

func (m *mockResponseWriter) Header() http.Header { return m.header }
func (m *mockResponseWriter) Write(p []byte) (int, error) {
	if !m.wroteHeader {
		m.WriteHeader(http.StatusOK)
	}
	return m.conn.Write(p)
}
func (m *mockResponseWriter) WriteHeader(code int) {
	if m.wroteHeader { return }
	m.wroteHeader = true
	fmt.Fprintf(m.conn, "HTTP/1.1 %d %s\r\n", code, http.StatusText(code))
	for k, vv := range m.header {
		for _, v := range vv {
			fmt.Fprintf(m.conn, "%s: %s\r\n", k, v)
		}
	}
	fmt.Fprint(m.conn, "\r\n")
}

func main() {
	port := flag.String("port", "9000", "API 서버 포트")
	useMock := flag.Bool("mock", false, "Mock 데이터 사용 여부")
	mysqlDSN := flag.String("dsn", "guss_user:1234@tcp(127.0.0.1:3306)/guss", "MySQL 연결 정보")
	maxConn := flag.Int("max_conn", 1000, "최대 동시 연결 수")
	flag.Parse()

	var repo repository.Repository
	var logRepo repository.LogRepository

	if *useMock {
		log.Println("--- [NOTICE] Mock 테스트 모드로 실행 중입니다 ---")
		repo = repository.NewMockRepository()
		logRepo = repository.NewMockLogRepository()
	} else {
		log.Println("--- [DATABASE] 로컬 MySQL에 연결을 시도합니다 ---")
		db, err := sql.Open("mysql", *mysqlDSN)
		if err != nil {
			log.Fatalf("DB 연결 초기화 실패: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			log.Fatalf("DB 연결 실패 (Ping): %v", err)
		}

		db.SetMaxOpenConns(*maxConn)
		repo = repository.NewMySQLRepository(db)
		logRepo = repository.NewMockLogRepository() 
	}

	congestionAlgo := &algo.RealTimeCalculator{}

	server := &api.Server{
		Repo:    repo,
		LogRepo: logRepo,
		Algo:    congestionAlgo,
	}

	limiter := tcp.NewConnLimiter(*maxConn)
	ln, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("포트 바인딩 실패: %v", err)
	}

	mux := http.NewServeMux()
	registerRoutes(mux, server)

	srv := &http.Server{
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Graceful Shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("서버 종료 중...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	log.Printf("GUSS API 서버 시작 (Port: %s, Mock: %v)", *port, *useMock)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		if limiter.Acquire() {
			tcp.AddConn()
			go func(c net.Conn) {
				defer limiter.Release()
				defer tcp.RemoveConn()
				defer c.Close()

				reader := bufio.NewReader(c)
				req, err := http.ReadRequest(reader)
				if err != nil { return }

				w := &mockResponseWriter{conn: c, header: make(http.Header)}
				mux.ServeHTTP(w, req)
			}(conn)
		} else {
			conn.Close()
		}
	}
}

func registerRoutes(mux *http.ServeMux, s *api.Server) {
	mux.HandleFunc("/api/register", s.HandleRegister)
	mux.HandleFunc("/api/login", s.HandleLogin)
	mux.HandleFunc("/api/gyms", s.HandleGetGyms)
	mux.HandleFunc("/api/gyms/", s.HandleGetGymDetail)
 
	mux.Handle("/api/reserve", s.AuthMiddleware(http.HandlerFunc(s.HandleReserve)))

	// Admin API (인증 + 관리자 권한 미들웨어)
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/dashboard", s.HandleGetAdminDashboard)
	adminMux.HandleFunc("/equip", s.HandleAddEquipment)
	mux.Handle("/api/admin/", s.AuthMiddleware(s.AdminMiddleware(adminMux)))
}
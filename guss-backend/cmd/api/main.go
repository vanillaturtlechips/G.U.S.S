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

	// MySQL 드라이버 로드
	_ "github.com/go-sql-driver/mysql"

	"guss-backend/internal/algo"
	"guss-backend/internal/api"
	"guss-backend/internal/repository"
	"guss-backend/pkg/tcp"
)

// corsMiddleware: 브라우저의 CORS 차단 방지
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// mockResponseWriter: TCP 연결을 통해 HTTP 응답을 전송하기 위한 래퍼
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
	if m.wroteHeader {
		return
	}
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
	// 1. 실행 옵션 설정 (기본 포트 9000)
	port := flag.String("port", "9000", "API 서버 포트")
	useMock := flag.Bool("mock", false, "Mock 데이터 사용 여부")
	mysqlDSN := flag.String("dsn", "guss_user:1234@tcp(127.0.0.1:3306)/guss", "MySQL 연결 정보")
	maxConn := flag.Int("max_conn", 1000, "최대 동시 연결 수")
	flag.Parse()

	var repo repository.Repository
	var logRepo repository.LogRepository

	// 2. 리포지토리 초기화
	if *useMock {
		log.Println("--- [NOTICE] Mock 테스트 모드로 실행 중입니다 ---")
		repo = repository.NewMockRepository()
		logRepo = repository.NewMockLogRepository()
	} else {
		log.Println("--- [DATABASE] MySQL 연결 시도 중... ---")
		db, err := sql.Open("mysql", *mysqlDSN)
		if err != nil {
			log.Fatalf("DB 초기화 실패: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			log.Fatalf("DB 연결 실패: %v", err)
		}

		db.SetMaxOpenConns(*maxConn)
		repo = repository.NewMySQLRepository(db)
		logRepo = repository.NewMockLogRepository()
	}

	// 3. 서버 인스턴스 생성
	server := &api.Server{
		Repo:    repo,
		LogRepo: logRepo,
		Algo:    &algo.RealTimeCalculator{},
	}

	mux := http.NewServeMux()
	registerRoutes(mux, server)

	srv := &http.Server{
		Handler:      corsMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// 4. Graceful Shutdown 처리
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("--- [SERVER] 종료 신호 감지 ---")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	// 5. TCP 서버 및 HTTP 핸들러 연결
	limiter := tcp.NewConnLimiter(*maxConn)
	ln, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("포트 바인딩 실패: %v", err)
	}

	log.Printf("GUSS API 서버 가동 중 (Port: %s)", *port)

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
				if err != nil {
					return
				}

				w := &mockResponseWriter{conn: c, header: make(http.Header)}
				srv.Handler.ServeHTTP(w, req)
			}(conn)
		} else {
			conn.Close()
		}
	}
}

// registerRoutes: 모든 API 경로를 핸들러와 연결 (404 방지)
func registerRoutes(mux *http.ServeMux, s *api.Server) {
	// [기존 사용자 경로]
	mux.HandleFunc("/api/register", s.HandleRegister)
	mux.HandleFunc("/api/login", s.HandleLogin)
	mux.HandleFunc("/api/gyms", s.HandleGetGyms)
	mux.HandleFunc("/api/gyms/", s.HandleGetGymDetail)
	mux.Handle("/api/reserve", s.AuthMiddleware(http.HandlerFunc(s.HandleReserve)))

	// [관리자 대시보드 경로 - admin.tsx 연동용]
	mux.HandleFunc("/api/dashboard", s.HandleDashboard)

	// [기구 관리 - 조회(GET) 및 등록(POST)]
	mux.HandleFunc("/api/equipments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.HandleAddEquipment(w, r)
		} else {
			s.HandleGetEquipments(w, r)
		}
	})

	// [기구 삭제 - /api/equipments/{id}]
	// 실제 구현 시 URL 파라미터 추출 로직이 필요할 수 있습니다.
	mux.HandleFunc("/api/equipments/", s.HandleDeleteEquipment)

	// [로그 조회 경로]
	mux.HandleFunc("/api/reservations", s.HandleGetReservations)
	mux.HandleFunc("/api/sales", s.HandleGetSales)
}

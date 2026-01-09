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

// corsMiddleware: 요청한 Origin을 동적으로 허용하여 브라우저의 CORS 차단을 방지합니다.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 브라우저의 사전 확인(Preflight) 요청 즉시 응답
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// mockResponseWriter: TCP 연결을 통해 직접 HTTP 응답을 작성하기 위한 래퍼입니다.
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
	// 1. 실행 옵션 설정 (Flags)
	port := flag.String("port", "9000", "API 서버 포트")
	useMock := flag.Bool("mock", false, "Mock 데이터 사용 여부")
	mysqlDSN := flag.String("dsn", "guss_user:1234@tcp(127.0.0.1:3306)/guss", "MySQL 연결 정보")
	maxConn := flag.Int("max_conn", 1000, "최대 동시 연결 수")
	flag.Parse()

	var repo repository.Repository
	var logRepo repository.LogRepository

	// 2. 데이터베이스 및 레포지토리 초기화
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

		// 연결 확인 (Ping)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			log.Fatalf("DB 연결 실패: %v", err)
		}

		db.SetMaxOpenConns(*maxConn)
		repo = repository.NewMySQLRepository(db)
		logRepo = repository.NewMockLogRepository() // 필요시 DynamoDB 등으로 교체 가능
	}

	// 3. API 서버 인스턴스 생성 (의존성 주입)
	// [중요] Panic 방지를 위해 Algo 엔진(&algo.RealTimeCalculator{})을 명시적으로 주입합니다.
	server := &api.Server{
		Repo:    repo,
		LogRepo: logRepo,
		Algo:    &algo.RealTimeCalculator{}, 
	}

	// 4. 라우팅 및 서버 설정
	mux := http.NewServeMux()
	registerRoutes(mux, server)

	// 미들웨어(CORS) 적용
	srv := &http.Server{
		Handler:      corsMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// 5. Graceful Shutdown (서버 안전 종료) 설정
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		log.Println("--- [SERVER] 종료 신호 감지, 시스템을 중단합니다 ---")
		
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("종료 중 오류 발생: %v", err)
		}
	}()

	// 6. TCP 리스너 및 연결 제한(Limiter) 실행
	limiter := tcp.NewConnLimiter(*maxConn)
	ln, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("포트 바인딩 실패: %v", err)
	}

	log.Printf("GUSS API 서버 가동 중 (Port: %s, MaxConn: %d)", *port, *maxConn)

	for {
		conn, err := ln.Accept()
		if err != nil {
			// 리스너가 닫힌 경우 루프 탈출
			select {
			case <-context.Background().Done():
				return
			default:
				continue
			}
		}

		// 동시 연결 수 제한 체크
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

				// HTTP 핸들러 실행
				w := &mockResponseWriter{conn: c, header: make(http.Header)}
				srv.Handler.ServeHTTP(w, req)
			}(conn)
		} else {
			// 접속 한도 초과 시 연결 즉시 종료
			conn.Close()
		}
	}
}

// registerRoutes: API 엔드포인트를 등록합니다.
func registerRoutes(mux *http.ServeMux, s *api.Server) {
	mux.HandleFunc("/api/register", s.HandleRegister)
	mux.HandleFunc("/api/login", s.HandleLogin)
	mux.HandleFunc("/api/gyms", s.HandleGetGyms)
	// 슬래시(/)를 붙여 하위 경로(/api/gyms/1 등)를 허용합니다.
	mux.HandleFunc("/api/gyms/", s.HandleGetGymDetail) 
 
	// 인증이 필요한 API
	mux.Handle("/api/reserve", s.AuthMiddleware(http.HandlerFunc(s.HandleReserve)))
}
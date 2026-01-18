package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	_ "github.com/go-sql-driver/mysql"

	"guss-backend/internal/algo"
	"guss-backend/internal/api"
	"guss-backend/internal/repository"
)

func main() {
	// 1. 실행 옵션 설정 (포트, DB 주소, SQS URL, DynamoDB 테이블)
	port := flag.String("port", "9000", "API 서버 포트")
	mysqlDSN := flag.String("dsn", "", "MySQL 연결 정보 (필수)")
	sqsURL := flag.String("sqs_url", "", "SQS FIFO 큐 주소 (필수)")
	dynamoTable := flag.String("dynamo_table", "GUSS-DEV-DDB", "DynamoDB 테이블 이름")
	flag.Parse()

	if *mysqlDSN == "" || *sqsURL == "" {
		log.Fatal("에러: -dsn과 -sqs_url 옵션은 필수입니다.")
	}

	// 2. AWS 설정 로드 및 클라이언트(SQS, DynamoDB) 초기화
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		log.Fatalf("AWS 설정 로드 실패: %v", err)
	}
	sqsClient := sqs.NewFromConfig(cfg)
	dynamoClient := dynamodb.NewFromConfig(cfg)

	// 3. MySQL 데이터베이스 연결
	db, err := sql.Open("mysql", *mysqlDSN)
	if err != nil {
		log.Fatalf("DB 초기화 실패: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	// 4. 레포지토리 및 서버 인스턴스 생성
	// Server 구조체에 DynamoClient와 DynamoTable을 주입합니다.
	repo := repository.NewMySQLRepository(db)
	server := &api.Server{
		Repo:         repo,
		Algo:         &algo.RealTimeCalculator{},
		SQSClient:    sqsClient,
		SQSURL:       *sqsURL,
		DynamoClient: dynamoClient,
		DynamoTable:  *dynamoTable,
	}

	// 5. 라우팅 설정
	mux := http.NewServeMux()

	// [공통 API]
	mux.HandleFunc("/api/register", server.HandleRegister)
	mux.HandleFunc("/api/login", server.HandleLogin)
	mux.HandleFunc("/api/dashboard", server.HandleDashboard)
	mux.HandleFunc("/api/gyms", server.HandleGetGyms)
	mux.HandleFunc("/api/gyms/", server.HandleGetGymDetail)

	// [체크인 API] - QR 스캔 시 호출되는 경로 (새로 추가됨)
	mux.HandleFunc("/api/checkin", server.HandleCheckIn)

	// [예약 API] - 인증 미들웨어 적용
	mux.Handle("/api/reserve", server.AuthMiddleware(http.HandlerFunc(server.HandleReserve)))
	mux.Handle("/api/reserve/cancel", server.AuthMiddleware(http.HandlerFunc(server.HandleCancelReservation)))

	// [관리자 API]
	mux.HandleFunc("/api/admin/sales", server.HandleGetSales)
	mux.HandleFunc("/api/admin/equipments", server.HandleGetEquipments)

	

	log.Printf("GUSS API 서버 가동 중 (Port: %s)", *port)
	if err := http.ListenAndServe(":"+*port, mux); err != nil {
		log.Fatalf("서버 실행 실패: %v", err)
	}
}

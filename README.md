# G.U.S.S (Gym User Support Service)

> 실시간 체육관 혼잡도 확인 및 예약 시스템  
> AWS 클라우드 인프라 기반 마이크로서비스 아키텍처

## 프로젝트 개요

G.U.S.S는 체육관 이용자에게 실시간 혼잡도 정보를 제공하고, 예약 및 QR 체크인 기능을 통해 편리한 이용 경험을 제공하는 클라우드 네이티브 서비스입니다.

### 주요 기능

- **실시간 혼잡도 모니터링**: 체육관별 현재 이용 인원 및 혼잡도 시각화
- **예약 시스템**: JWT 인증 기반 30분 단위 시간대 예약
- **QR 체크인**: 예약 확인 및 현장 입장 처리
- **FCM 푸시 알림**: 예약 확정 및 혼잡도 알림
- **관리자 대시보드**: 예약 현황, 매출 통계, 기구 관리

---

## 아키텍처

### 기술 스택

**Frontend**
- React 19 + TypeScript
- Vite (빌드 도구)
- TailwindCSS 4 + Framer Motion
- React Router v7
- Chart.js (통계 시각화)
- Firebase FCM (푸시 알림)
- QRCode.react

**Backend**
- Go 1.x (멀티 커맨드 구조)
- Gorilla Mux / net/http
- JWT 인증 (golang-jwt)
- Bcrypt 패스워드 해싱

**Infrastructure**
- AWS EKS (Kubernetes 오케스트레이션)
- AWS Lambda (비동기 워커)
- Amazon RDS MySQL 8.0
- Amazon DynamoDB (예약 로그)
- Amazon SQS FIFO (메시지 큐)
- Amazon SNS (알림)
- AWS Secrets Manager (민감 정보 관리)
- Terraform (IaC)

### 시스템 구성도

```
[사용자] → [ALB] → [Nginx (Public Subnet)]
                      ↓
              [Internal ALB] → [Go API Server (Private Subnet)]
                                  ↓
                    ┌─────────────┼─────────────┐
                    ↓             ↓             ↓
                [RDS MySQL]   [DynamoDB]    [SQS FIFO]
                                              ↓
                                        [Lambda Worker]
                                              ↓
                                        [Firebase FCM]
```

---

## 프로젝트 구조

```
G.U.S.S/
├── frontend/                    # React 프론트엔드
│   ├── src/
│   │   ├── pages/              # 페이지 컴포넌트
│   │   │   ├── Dashboard.tsx   # 메인 대시보드 (체육관 목록)
│   │   │   ├── guss.tsx        # 체육관 상세 (예약/QR)
│   │   │   ├── Login.tsx       # 로그인
│   │   │   ├── register.tsx    # 회원가입
│   │   │   └── admin.tsx       # 관리자 페이지
│   │   ├── components/
│   │   │   ├── reservation/    # 예약 관련 컴포넌트
│   │   │   └── charts/         # 차트 컴포넌트
│   │   ├── api/
│   │   │   └── axios.ts        # API 클라이언트 설정
│   │   └── firebase/
│   │       └── firebaseConfig.ts
│   └── public/
│       └── firebase-messaging-sw.js  # FCM Service Worker
│
├── guss-backend/               # Go 백엔드
│   ├── cmd/
│   │   ├── api/                # API 서버 진입점
│   │   │   └── main.go
│   │   ├── worker/             # Lambda 워커
│   │   │   └── main.go
│   │   └── lambda-worker/      # 대체 워커 구현
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers.go     # HTTP 핸들러
│   │   │   ├── middleware.go   # JWT 인증 미들웨어
│   │   │   └── swagger.yaml    # API 문서
│   │   ├── repository/
│   │   │   ├── mysql_repository.go   # MySQL 레포지토리
│   │   │   ├── dynamo_repo.go        # DynamoDB 레포지토리
│   │   │   └── mock_repo.go          # 테스트용 Mock
│   │   ├── domain/
│   │   │   └── models.go       # 도메인 모델 (User, Gym, Reservation 등)
│   │   ├── auth/
│   │   │   ├── jwt.go          # JWT 토큰 생성/검증
│   │   │   └── password.go     # Bcrypt 해싱
│   │   ├── algo/
│   │   │   └── congestion.go   # 혼잡도 계산 알고리즘
│   │   └── infrastructure/
│   │       └── aws/            # AWS SDK 래퍼
│   ├── pkg/
│   │   └── tcp/                # TCP 메트릭/헬퍼
│   └── db/
│       └── schema.sql          # MySQL 스키마 정의
│
├── terraform/                  # Infrastructure as Code
│   ├── main.tf                 # Provider 설정
│   ├── variables.tf            # 변수 정의
│   ├── network.tf              # VPC, Subnet, IGW, NAT
│   ├── compute.tf              # EC2, ALB, Target Group
│   ├── data_messaging.tf       # RDS, DynamoDB, SQS, SNS
│   ├── security.tf             # Security Group
│   ├── iam.tf                  # IAM Role & Policy
│   └── monitoring.tf           # CloudWatch
│
├── deploy_api.sh               # API 서버 배포 스크립트
├── deploy_worker.sh            # Lambda 워커 배포 스크립트
└── README.md
```

---

## 시작하기

### 사전 요구사항

- Node.js 18+
- Go 1.21+
- AWS CLI v2
- Terraform 1.5+
- MySQL 8.0 클라이언트

### 1. 환경 변수 설정

**Backend (.env 또는 실행 플래그)**
```bash
# MySQL 연결 정보
DSN="user:password@tcp(rds-endpoint:3306)/guss?parseTime=true"

# AWS 리소스
SQS_URL="https://sqs.ap-northeast-2.amazonaws.com/account-id/GUSS-PRD-SQS-RSV"
DYNAMO_TABLE="GUSS-PRD-DDB"

# JWT 시크릿
JWT_SECRET="your-secret-key"
```

**Frontend (.env)**
```bash
VITE_API_BASE_URL=http://localhost:9000
VITE_FIREBASE_API_KEY=your-firebase-api-key
VITE_FIREBASE_PROJECT_ID=your-project-id
```

### 2. 데이터베이스 초기화

```bash
mysql -h <RDS_ENDPOINT> -u admin -p < guss-backend/db/schema.sql
```

### 3. 백엔드 실행

```bash
cd guss-backend

# API 서버 빌드 및 실행
go build -o guss-api cmd/api/main.go
./guss-api -port 9000 -dsn "$DSN" -sqs_url "$SQS_URL"
```

### 4. 프론트엔드 실행

```bash
cd frontend
npm install
npm run dev
```

브라우저에서 `http://localhost:5173` 접속

---

## 배포

### API 서버 배포 (EC2)

```bash
./deploy_api.sh
```

- Go 바이너리 빌드
- systemd 서비스 재시작
- 상태 확인

### Lambda 워커 배포

```bash
./deploy_worker.sh
```

- GOOS=linux 크로스 컴파일
- bootstrap 바이너리 생성
- Lambda 함수 코드 업데이트
- SSM에서 Firebase 키 주입

### Terraform 인프라 프로비저닝

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

---

## API 엔드포인트

### 인증
- `POST /api/register` - 회원가입 (Bcrypt 해싱)
- `POST /api/login` - 로그인 (JWT 발급)

### 체육관
- `GET /api/dashboard` - 전체 체육관 목록 + 혼잡도
- `GET /api/gyms` - 체육관 목록
- `GET /api/gyms/{id}` - 체육관 상세 정보

### 예약 (JWT 필수)
- `POST /api/reserve` - 예약 생성
- `POST /api/reserve/cancel` - 예약 취소
- `GET /api/reserve/active` - 활성 예약 조회

### 체크인
- `POST /api/checkin` - QR 체크인

### 관리자 (JWT + Admin 권한)
- `GET /api/admin/reservations` - 예약 현황
- `GET /api/admin/sales` - 매출 통계
- `GET /api/admin/equipments` - 기구 목록
- `POST /api/admin/equipments` - 기구 추가
- `PUT /api/admin/equipments/{id}` - 기구 수정
- `DELETE /api/admin/equipments/{id}` - 기구 삭제

상세 API 문서: `guss-backend/internal/api/swagger.yaml`

---

## 보안

- **비밀번호**: Bcrypt 해싱 (cost 10)
- **인증**: JWT (HS256) + Bearer Token
- **네트워크**: Private Subnet 격리, Security Group 최소 권한
- **민감 정보**: AWS Secrets Manager / SSM Parameter Store

---

## 모니터링

- **CloudWatch Logs**: API 서버 및 Lambda 로그
- **CloudWatch Metrics**: ALB, RDS, DynamoDB 메트릭
- **SNS 알림**: 시스템 장애 알림

---

## 테스트

```bash
# Backend 테스트
cd guss-backend
go test ./...

# Frontend 테스트
cd frontend
npm run test
```

---

## 기여

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 라이선스

이 프로젝트는 MIT 라이선스를 따릅니다. 자세한 내용은 `LICENSE` 파일을 참조하세요.

---

## 팀

**2조 클라우드 엔지니어링 프로젝트**

- 프로젝트 기간: 2025.01 ~ 2025.04
- 기술 스택: Go, React, AWS, Terraform

---

## 문의

프로젝트 관련 문의사항은 이슈를 등록해주세요.

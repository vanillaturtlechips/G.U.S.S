#!/bin/bash

# 1. 동적 경로 설정 (어디서 실행하든 현재 폴더를 기준으로 잡음)
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
BINARY_NAME="guss-api"

echo "--- [INFO] 현재 작업 경로: $PROJECT_ROOT ---"
cd "$PROJECT_ROOT"

# 2. 로컬 변경 사항 처리 (Pull 충돌 방지)
echo "--- [GIT] 로컬 변경사항 Stash 및 Pull 시작 ---"
git stash
git pull origin myong/lambda

# 3. API 서버 빌드 (cmd/api/main.go -> guss-api)
echo "--- [BUILD] Go API 서버 빌드 중... ---"
go build -v -o $BINARY_NAME cmd/api/main.go

if [ $? -eq 0 ]; then
    echo "--- [SUCCESS] 빌드 성공: $BINARY_NAME ---"
else
    echo "--- [ERROR] 빌드 실패! 로그를 확인하세요. ---"
    exit 1
fi

# 4. 실행 권한 부여
chmod +x $BINARY_NAME

# 5. Systemd 서비스 재시작
echo "--- [DEPLOY] guss-api 서비스 재시작 중... ---"
sudo systemctl restart guss-api

# 6. 최종 상태 확인
echo "--- [STATUS] 서비스 상태 확인 ---"
systemctl status guss-api --no-pager

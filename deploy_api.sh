#!/bin/bash

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
BINARY_NAME="guss-api"

echo "--- [INFO] 현재 작업 경로: $PROJECT_ROOT ---"
cd "$PROJECT_ROOT"

# [중요] 스크립트 수정을 유지하기 위해 빌드 전 stash/pull은 잠시 주석 처리하거나 
# 빌드 후에 수행하도록 순서를 조정하는 것이 안전합니다.
# git stash
# git pull origin myong/lambda

# 1. backend 폴더로 이동 (go.mod가 있는 곳)
echo "--- [MOVE] guss-backend 폴더로 이동 ---"
cd guss-backend

# 2. API 서버 빌드
echo "--- [BUILD] Go API 서버 빌드 중... ---"
# 루트 폴더에 바이너리를 생성합니다.
go build -v -o ../$BINARY_NAME cmd/api/main.go

if [ $? -eq 0 ]; then
    echo "--- [SUCCESS] 빌드 성공: $BINARY_NAME ---"
else
    echo "--- [ERROR] 빌드 실패! 로그를 확인하세요. ---"
    exit 1
fi

cd "$PROJECT_ROOT"
chmod +x $BINARY_NAME

echo "--- [DEPLOY] guss-api 서비스 재시작 중... ---"
sudo systemctl restart guss-api

echo "--- [STATUS] 서비스 상태 확인 ---"
systemctl status guss-api --no-pager

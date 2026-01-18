#!/bin/bash

# 1. guss-backend 디렉토리로 이동
cd guss-backend

# 2. 빌드 실행 (API 서버용)
echo "Building GUSS API binary..."
# -o 옵션 이름을 guss-api로 차별화합니다.
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o guss-api cmd/api/main.go

# 빌드 성공 여부 확인
if [ $? -ne 0 ]; then
    echo "API Build failed!"
    exit 1
fi

echo "API Build finished: guss-api"

# 3. (선택 사항) 서버 재시작 로직
# 지금은 베스천에서 직접 테스트하시니까, 여기까지만 하고 나중에 
# AMI로 배포할 때 이 밑에 실행 코드를 넣으면 됩니다.

echo "Ready to run!"

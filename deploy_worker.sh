#!/bin/bash

# 1. 환경 변수 및 설정
FUNCTION_NAME="GUSS-Reservation-Worker"
SSM_PARAMETER_NAME="/guss/worker/firebase-key"
TABLE_NAME="GUSS-DEV-DDB"

# 2. guss-backend 디렉토리로 이동
cd guss-backend

# 3. 빌드 실행 (Amazon Linux용 64비트 바이너리)
echo "Building Go binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o bootstrap cmd/worker/main.go

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

# 4. 압축 (이제 JSON 파일 없이 bootstrap만 압축합니다)
echo "Packaging bootstrap..."
zip function.zip bootstrap

# 5. AWS Lambda 코드 업데이트
echo "Deploying code to AWS Lambda..."
aws lambda update-function-code \
    --function-name $FUNCTION_NAME \
    --zip-file fileb://function.zip

# 코드 업데이트가 반영될 때까지 잠시 대기 (안정적인 설정을 위해)
echo "Waiting for function update to complete..."
aws lambda wait function-updated --function-name $FUNCTION_NAME

# 6. SSM에서 Firebase 키 가져오기
echo "Retrieving Firebase secret from SSM..."
FIREBASE_CONFIG=$(aws ssm get-parameter --name "$SSM_PARAMETER_NAME" --with-decryption --query "Parameter.Value" --output text)

if [ -z "$FIREBASE_CONFIG" ]; then
    echo "Failed to retrieve Firebase config from SSM!"
    exit 1
fi

# 7. Lambda 환경 변수 업데이트 (SSM에서 가져온 JSON 주입)
echo "Updating Lambda environment variables..."
aws lambda update-function-configuration \
    --function-name $FUNCTION_NAME \
    --environment "Variables={TABLE_NAME=$TABLE_NAME,FIREBASE_CONFIG='$FIREBASE_CONFIG'}"

# 8. 정리
rm bootstrap function.zip

echo "Deployment finished successfully!"

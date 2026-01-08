module guss-backend

go 1.24.0

toolchain go1.24.11

require (
	github.com/aws/aws-sdk-go-v2 v1.41.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.29
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.53.5
	github.com/golang-jwt/jwt/v5 v5.3.0
	golang.org/x/crypto v0.46.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.16 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.32.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.16 // indirect
	github.com/aws/smithy-go v1.24.0 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
)

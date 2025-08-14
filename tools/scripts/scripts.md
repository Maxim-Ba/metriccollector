go test --coverprofile=coverage.out ./... && go tool cover -func=coverage.out - общее покрытие
gofmt -l .
gofmt -w .


./cmd/multichecker/multichecker ./...


//---- генерация ключей 

openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -outform PEM -pubout -out public.pem



//---- генерация прото


protoc --go_out=. --go_opt=paths=source_relative   --go-grpc_out=. --go-grpc_opt=paths=source_relative   proto/metrics.proto

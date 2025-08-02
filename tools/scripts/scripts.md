go test --coverprofile=coverage.out ./... && go tool cover -func=coverage.out - общее покрытие
gofmt -l .
gofmt -w .


./cmd/multichecker/multichecker ./...


//---- генерация ключей 

openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -outform PEM -pubout -out public.pem

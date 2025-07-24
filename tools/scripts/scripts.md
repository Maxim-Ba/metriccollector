go test --coverprofile=coverage.out ./... && go tool cover -func=coverage.out - общее покрытие
gofmt -l .
gofmt -w .

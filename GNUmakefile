
tidy:
	@go mod tidy

upd:
	@go get -u ./...

lint:
	@golangci-lint run pkg/...

lintfix:
	@golangci-lint run pkg/... --fix

test:
	@go test -v ./...
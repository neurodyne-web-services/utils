
tidy:
	@go mod tidy --compat=1.19

upd:
	@./upd.sh

lint:
	@golangci-lint run pkg/...

lintfix:
	@golangci-lint run pkg/... --fix

test:
	@go test -v ./...
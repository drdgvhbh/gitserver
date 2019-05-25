start:
	go run ./cmd/main.go

test:
	go test ./internal/...

test-e2e:
	go test ./test
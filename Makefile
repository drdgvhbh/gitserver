start:
	go run ./cmd/main.go

test-unit:
	go test ./internal/...

test-e2e:
	go test ./test
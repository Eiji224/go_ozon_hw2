.PHONY: version
version:
	go version
	protoc --version

.PHONY: run-all
run-all:
	cd cart && go run ./cmd/api/main.go
.PHONY: version
version:
	go version
	protoc --version

.PHONY: run-all
run-all:
	cd cart && go run ./cmd/api/main.go

.PHONY: generate
generate:
	cd product_service && protoc -I ./ \
	--go_out ./pkg --go_opt paths=source_relative \
	--go-grpc_out ./pkg --go-grpc_opt paths=source_relative \
	api/product_service.proto
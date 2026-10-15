.PHONY: version
version:
	go version
	protoc --version

.PHONY: run-all
run-all:
	trap 'kill 0' EXIT INT TERM; \
	(cd product_service && go run ./cmd/product_service/main.go) & \
	(cd cart && go run ./cmd/api/main.go) & \
	wait 

.PHONY: generate
generate:
	cd product_service && protoc -I ./ \
	--go_out ./pkg --go_opt paths=source_relative \
	--go-grpc_out ./pkg --go-grpc_opt paths=source_relative \
	api/product_service.proto
BIN := "./bin/app"
DOCKER_IMG="xlabs-app"

build:
	go build -v -o $(BIN) ./cmd

run: build
	$(BIN) -config ./configs/config.yaml -admin_name admin -admin_password admin

docker-build:
	docker build \
		-t $(DOCKER_IMG) \
        -f Dockerfile .
docker-run: docker-build
	docker run  \
		-p 50051:50051 \
		$(DOCKER_IMG)

test:
	go test -cover -v ./...

generate:
	rm -rf internal/api/grpc/pb
	mkdir -p internal/api/grpc/pb
	protoc --proto_path=api/ --go_out=internal/api/grpc/pb	--go-grpc_out=internal/api/grpc/pb api/*.proto

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run docker-build docker-run test lint






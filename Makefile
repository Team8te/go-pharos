
.PHONY: generate
generate:
	$(info Start generation...)
	go install github.com/swaggo/swag/cmd/swag@latest
	cd ./cmd/go-pharos/ ; pwd ; swag init ; mv ./docs ../../

.PHONY: generate-mocks
generate-mocks:
	$(info Start mock generation...)

.PHONY: docker-build
docker-build:
	docker build -f ./docker/Dockerfile -t go-pharos-develop .

.PHONY: install
install:
	go build -v -o /usr/local/bin/ ./cmd/go-pharos/...
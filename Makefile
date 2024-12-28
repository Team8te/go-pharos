

generate:
	$(info Start generation...)
	go install github.com/swaggo/swag/cmd/swag@latest
	cd ./cmd/go-pharos/ ; pwd ; swag init ; mv ./docs ../../


generate-mocks:
	$(info Start mock generation...)
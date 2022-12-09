SHELL=/bin/bash

.PHONY: all build deploy clean dev proto build-client build-server deploy-client build-server

include .env

PATH_CURRENT := $(shell pwd)
PATH_BUILT := $(PATH_CURRENT)/build/
GIT_COMMIT_LOG := $(shell git log --oneline -1 HEAD)


build-gateway:
	echo "current commit: ${GIT_COMMIT_LOG}"
	go mod tidy
	env GOOS=linux GOARCH=amd64 go build -v -o ./build/mpc_gateway -ldflags "-X 'main.GitCommitLog=${GIT_COMMIT_LOG}'" ./cmd/mpc_gateway

build-client:
	echo "current commit: ${GIT_COMMIT_LOG}"
	go mod tidy
	env GOOS=linux GOARCH=amd64 go build -v -o ./build/mpc_client -ldflags "-X 'main.GitCommitLog=${GIT_COMMIT_LOG}'" ./cmd/mpc_client

deploy-gateway: clean build-gateway
	gcloud run deploy --source . --region asia-southeast1 --project ${GCP_PROJECT}; \

deploy-client: clean build-client
	gcloud run deploy --source . --region asia-southeast1 --project ${GCP_PROJECT}; \

clean:
	rm -fr "${PATH_BUILT}/mpc_client"; \
	rm -fr "${PATH_BUILT}/mpc_gateway"; \

build:
	go build -v -o ./build/server-local -ldflags "-X 'main.GitCommitLog=${GIT_COMMIT_LOG}'"

migrate-client-up:
	migrate -path db/client/migration -database "$(DB_URL)" -verbose up

migrate-client-down:
	migrate -path db/client/migration -database "$(DB_URL)" -verbose down

migrate-gateway-up:
	migrate -path db/gateway/migration -database "$(DB_URL)" -verbose up

migrate-gateway-down:
	migrate -path db/gateway/migration -database "$(DB_URL)" -verbose down

sqlc: sqlc-client sqlc-gateway
sqlc-client:
	sqlc generate --file sqlc-client.yaml
sqlc-gateway:
	sqlc generate --file sqlc-gateway.yaml

gateway: sqlc-gateway
	go run ./cmd/mpc_gateway/main.go

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--descriptor_set_out descriptor.pb \
	proto/*.proto

SHELL=/bin/bash

.PHONY: all build deploy clean dev build-linux proto

include .env

PATH_CURRENT := $(shell pwd)
PATH_BUILT := $(PATH_CURRENT)/build/server
GIT_COMMIT_LOG := $(shell git log --oneline -1 HEAD)

all: build-linux deploy clean

build-linux:
	echo "current commit: ${GIT_COMMIT_LOG}"
	go mod tidy
	env GOOS=linux GOARCH=amd64 go build -v -o ./build/server -ldflags "-X 'main.GitCommitLog=${GIT_COMMIT_LOG}'"

deploy: clean build-linux
	gcloud run deploy --source . --region asia-southeast1 --project ${GCP_PROJECT}; \

clean:
	rm -fr "${PATH_BUILT}"; \
	echo "Clean built."

build:
	go build -v -o ./build/server-local -ldflags "-X 'main.GitCommitLog=${GIT_COMMIT_LOG}'"

dev: build
	./build/server-local

server:
	go run main.go

migrate-up:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--descriptor_set_out descriptor.pb \
	proto/*.proto

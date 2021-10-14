
deps:
	go mod download
	go get -u google.golang.org/grpc
	go get -u github.com/golang/protobuf/protoc-gen-go
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

build_c: 
	go build -o spqr-c cmd/coordinator/main.go

build_proxy: 
	go build -o spqr-rr cmd/router/main.go

build_world: 
	go build -o spqr-world cmd/world/main.go

build_worldmock:
	go build -o spqr-worldmock ./test/worldmock

build: build_c build_proxy build_world build_worldmock

gen: gogen yaccgen

gogen:
	protoc --go_out=./router --go_opt=paths=source_relative --go-grpc_out=./router --go-grpc_opt=paths=source_relative \
	protos/* 

init:
	 go mod download
	 go mod vendor

build_images:
	docker-compose build spqrbase shardbase

test: build_images
	docker-compose up --remove-orphans --exit-code-from client --build router shard1 shard2 client

run: build_images
	docker-compose up -d --remove-orphans --build router coordinator world1 shard1 shard2
	docker-compose build client
	docker-compose run --entrypoint /bin/bash client

yaccgen:
	goyacc -o yacc/console/sql.go -p yy yacc/console/sql.y

.PHONY: build gen

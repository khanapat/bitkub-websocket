PROJECT?=bitkub-websocket
RELEASE?=$(shell git tag --points-at HEAD)

GOOS?=linux
GOARCH?=amd64

APP?=goapp
PORT?=9090

CACHE_IMAGE?=$$(docker images --filter "dangling=true" -q --no-trunc)

run:
	go run main.go

clean:
	rm -f $(APP)

test: clean
	go test -v -cover ./...

container: test
	docker build . --no-cache -t $(PROJECT):$(RELEASE) -f build/Dockerfile
	docker rmi $(CACHE_IMAGE)

create-app:
	docker-compose -f ./build/docker-compose.yaml up -d

delete-app:
	docker-compose -f ./build/docker-compose.yaml down
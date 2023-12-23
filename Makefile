PROJECT_NAME=go-spider-be

export GOPATH :=$(HOME)/go
export PATH :=$(GOPATH)/bin:$(PATH)



# command controll
test:
	go test -cover --race -v -failfast ./...

# set up project
up:
	docker-compose up -d
down:
	docker-compose down

shell:
	docker exec -it go-spider-be bash
	
mongo-shell:
	docker exec -it mongodb bash
	
redis-shell:
	docker exec -it redis bash

generate:
	go generate ./...

build:
	env GOOS=linux GOARCH=amd64 go build -o go-spider-be

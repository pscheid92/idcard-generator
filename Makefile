NOW := $(shell date +"%Y%m%d-%H%M%S")

run:
	go run main.go

build:
	GOOS=linux go build .
	docker build  . \
		-t registry.gitlab.com/pscheid92/idcardgenerator:$(NOW) \
		-t registry.gitlab.com/pscheid92/idcardgenerator:latest

publish: build
	docker push registry.gitlab.com/pscheid92/idcardgenerator:$(NOW)
	docker push registry.gitlab.com/pscheid92/idcardgenerator:latest



all:
	rm -rf release
	mkdir release
	GOOS=linux GOARCH=amd64 go build -o ./release/ioa-httpServer ./cmd/httpServer/main.go
	GOOS=linux GOARCH=amd64 go build -o ./release/ioa-proxy ./cmd/proxy/main.go
	GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./plugins/cache.so ./plugins/cache.go
	GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./plugins/size.so ./plugins/size.go
	GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./plugins/rate.so ./plugins/rate.go
	GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./plugins/black.so ./plugins/black.go
	GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./plugins/white.so ./plugins/white.go
	GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./plugins/cors.so ./plugins/cors.go
	GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./plugins/jwt.so ./plugins/jwt.go
	GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./plugins/copy_request.so ./plugins/copy_request.g
	GOOS=linux GOARCH=amd64	go build -buildmode=plugin -o ./plugins/default_response.so ./plugins/default_response.go
	GOOS=linux GOARCH=amd64	go build -buildmode=plugin -o ./plugins/token_to_userId.so ./plugins/token_to_userId.go
	#rm -rf plugins/*.so
	#@for FILE in $(shell ls plugins); do \
	#	echo "building " $$FILE ;\
	#	echo ">>> $$(basename $$FILE .go)";\
	#	BASENAME=$$(basename $$FILE .go);\
	#	echo $(BASENAME);\
	#	go build -buildmode=plugin -o ./plugins/$$BASENAME.so ./plugins/size.go ;\
	#done
dev:
	go build -buildmode=plugin -o ./plugins/token_to_userId.so ./plugins/token_to_userId.go
	go build -buildmode=plugin -o ./plugins/default_response.so ./plugins/default_response.go
	go build -buildmode=plugin -o ./plugins/cache.so ./plugins/cache.go
	go build -buildmode=plugin -o ./plugins/copy_request.so ./plugins/copy_request.go
	go build -buildmode=plugin -o ./plugins/size.so ./plugins/size.go
	go build -buildmode=plugin -o ./plugins/rate.so ./plugins/rate.go
	go build -buildmode=plugin -o ./plugins/black.so ./plugins/black.go
	go build -buildmode=plugin -o ./plugins/white.so ./plugins/white.go
	go build -buildmode=plugin -o ./plugins/cors.so ./plugins/cors.go
	go build -buildmode=plugin -o ./plugins/jwt.so ./plugins/jwt.go
	go run cmd/proxy/main.go -config ./

linux:
	docker run -v "$$GOPATH":/go --rm -v "$$PWD":/go/src/myapp -w /go/src/myapp -e=GOOS=linux -e=GOARCH=amd64 -e=GO111MODULE=on  golang:latest make
clean:
	@rm -rf ./release/*
gotool:
	gofmt -w .cc
	go tool vet . |& grep -v vendor;true
help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make gotool - run go tool 'fmt' and 'vet'"
	@echo "make ca - generate ca files"

.PHONY: clean gotool ca help...

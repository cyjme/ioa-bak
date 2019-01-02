all:
	GOOS=linux GOARCH=amd64 go build -o ./release/ioa ./cli/main.go
	GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./plugins/size.so ./plugins/size.go
	GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./plugins/rate.so ./plugins/rate.go
	#rm -rf plugins/*.so
	#@for FILE in $(shell ls plugins); do \
	#	echo "building " $$FILE ;\
	#	echo ">>> $$(basename $$FILE .go)";\
	#	BASENAME=$$(basename $$FILE .go);\
	#	echo $(BASENAME);\
	#	go build -buildmode=plugin -o ./plugins/$$BASENAME.so ./plugins/size.go ;\
	#done
dev:
	go build -buildmode=plugin -o ./plugins/size.so ./plugins/size.go
	go build -buildmode=plugin -o ./plugins/rate.so ./plugins/rate.go
	go run cli/main.go
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

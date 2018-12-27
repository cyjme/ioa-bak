all:
	GOOS=linux GOARCH=amd64 go build -o ./release/ioa ./cli/main.go
	GOOS=linux GOARCH=amd64 go build -buildmode=plugin -o ./plugins/size.so ./plugins/size.go
clean:
	@rm -rf ./release/*
gotool:
	gofmt -w .
	go tool vet . |& grep -v vendor;true
help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make gotool - run go tool 'fmt' and 'vet'"
	@echo "make ca - generate ca files"

.PHONY: clean gotool ca help...

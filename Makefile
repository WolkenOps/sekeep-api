.PHONY: deps clean build

deps:
	dep ensure

clean:
	rm -rf ./create/create
	rm -rf ./read/read
	rm -rf ./delete/delete
	rm -rf ./update/update
	rm -rf ./auth/auth

build:
	GOOS=linux GOARCH=amd64 go build -o create/create ./create
	GOOS=linux GOARCH=amd64 go build -o read/read ./read
	GOOS=linux GOARCH=amd64 go build -o delete/delete ./delete
	GOOS=linux GOARCH=amd64 go build -o update/update ./update
	GOOS=linux GOARCH=amd64 go build -o auth/auth ./auth
.PHONY: deps clean build

deps:
	dep ensure

clean:
	rm -r target/

build:
	mkdir -p target/
	GOOS=linux GOARCH=amd64 go build -o target/create ./api/create
	GOOS=linux GOARCH=amd64 go build -o target/read ./api/read
	GOOS=linux GOARCH=amd64 go build -o target/delete ./api/delete
	GOOS=linux GOARCH=amd64 go build -o target/update ./api/update
	GOOS=linux GOARCH=amd64 go build -o target/auth ./api/auth
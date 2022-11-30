
ifndef version
	version := 1.0.0-alpha
endif

.PHONY: build

ifndef $(spec)
	echo "spec location needs to be specified"
	exit 1
endif

get-release:
	@git describe --tags `git rev-list --tags --max-count=1`
	
build:
	go build -ldflags="-X main.BuildVersion=$(version)" && \
	go build -o bin/boxee 

install-cli: build
	mv bin/boxee ~/.local/bin

generate-code:
	oapi-codegen -old-config-style --generate types,client -package main $(spec) > oapi_client.gen.go


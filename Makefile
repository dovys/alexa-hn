SHELL = /bin/bash
MAKEFLAGS+=-s

build:
	GOOS=linux go build -o build/svc cmd/*

deploy: build
	gcloud app deploy

.PHONY: build deploy
.PHONY: build 

build:
	go build -o main

run: build
	./main
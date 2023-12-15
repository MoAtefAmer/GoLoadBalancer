.PHONY: build 

build:
	go build -o main

run: build
	./main


pyserver:;python3 -m http.server ${p} --directory server${p}
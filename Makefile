all: build

.PHONY: build
build:
	mkdir -p bin
	go build -o bin/gocsv main.go

.PHONY: install
install:
	go install ./main.go

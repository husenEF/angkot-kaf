.PHONY: run build clean

# Load environment variables from .env file
include .env
export

run:
	go run main.go

build:
	go build -o angkot-kaf main.go

clean:
	rm -f angkot-kaf
	rm -f angkot-kaf.exe
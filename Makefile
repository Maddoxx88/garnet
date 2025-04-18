run:
	go run main.go

build:
	go build -o garnet main.go

test:
	go test ./store

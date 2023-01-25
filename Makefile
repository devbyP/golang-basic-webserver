build:
	@go build -o ./bin/web main.go
run: build
	./bin/web

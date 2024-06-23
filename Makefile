build:
	@go build -o bin/memtravel memtravel.go

test:
	@go test -v ./...
	
run: build
	@./bin/memtravel
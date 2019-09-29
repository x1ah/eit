build:
	@go build -o ./out/eit cmd/eit.go

run: build
	@./out/eit

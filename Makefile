build:
	@go build -o ./out/eit cmd/eit/cli.go

run: build
	@./out/eit

test:
	echo OK

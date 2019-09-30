build:
	@go build -o ./out/eit cmd/eit/eit.go

run: build
	@./out/eit

test:
	echo OK

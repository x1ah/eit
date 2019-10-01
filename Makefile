build:
	@go build -race -o ./out/eit .

run: build
	@./out/eit

test:
	echo OK

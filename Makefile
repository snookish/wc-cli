.PHONY: build run test bench clean

DIST := dist
BINARY := counter

build:
	@mkdir -p $(DIST)
	@go build -o $(DIST)/$(BINARY) .

run: build
	@./$(DIST)/$(BINARY) $(ARGS)

test:
	@go test -v ./...

bench:
	@go test -bench=. -benchmem ./...

clean:
	@rm -rf $(DIST)
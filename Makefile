# Makefile

# Define the binary name
BINARY_NAME=myapp

# Build the binary
build:
	go build -o $(BINARY_NAME) main.go

# Run the binary
run: build
	./$(BINARY_NAME)

# Watch for changes and restart the server
watch:
	nodemon --exec "make run" --ext go

run:
	@echo "Starting the application..."
	go run main.go

test:
	@echo "Running tests..."
	go test ./...

build:
	@echo "Building the application..."
	go build -o bin/products_app main.go
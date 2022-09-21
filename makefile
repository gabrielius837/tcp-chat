BINARY_NAME=server
OUTPUT=bin
 
build:
	@echo "Building..."
	go build -o ${OUTPUT}/${BINARY_NAME} server.go
 
run: build
	@echo "Running..."
	./${OUTPUT}/${BINARY_NAME}
 
clean:
	@echo "Cleaning..."
	go clean
	rm -rfd ./${OUTPUT}
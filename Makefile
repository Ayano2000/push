# VARIABLES
build_dir := ./bin
source := ./cmd/push/main.go

.PHONY: build
build: clean
	@echo "Building the application..."
	go build -o $(build_dir)/push $(source)

.PHONY: run
run: build
	@echo "Running the application with argument: $(env)"
	$(build_dir)/push $(env)

.PHONY: clean
clean:
	rm -rf $(build_dir)

.PHONY: down
down:
	docker-compose down -v

.PHONY: up
up:
	docker-compose up -d

.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  make build  						   - Build the application"
	@echo "  make run env=<development|production> - Build and run the application with a target environment"
	@echo "  make clean                            - Remove build artifacts"
	@echo "  make help                             - Show this help message"
	@echo "  make up                               - Start the projects containers in the background"
	@echo "  make down                             - Stop the projects containers"


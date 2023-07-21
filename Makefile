PROJECT_NAME=furrybot
BUILD_DIR=build

settings:
	cp settings.example.json settings.dev.json
	cp settings.example.json settings.json
	@echo "----------------------------------------------"
	@echo "Created settings based on example"
	@echo "Please modify them before running the program"
	@echo "\"settings.dev.json\" for local development"
	@echo "\"settings.json\" for use in docker container"
	@echo "----------------------------------------------"

run: export FURRYBOT_CONFIG_FILE=settings.dev.json
run: 
	go run main.go

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(PROJECT_NAME) main.go

docker-build:
	docker compose -f docker-compose.yaml build

up:
	docker compose -f docker-compose.yaml up -d

down:
	docker compose -f docker-compose.yaml down

logs:
	docker compose -f docker-compose.yaml logs -f

clean:
	rmdir $(BUILD_DIR)

.PHONY: build run docker-build
PROJECT_NAME=furrybot
BUILD_DIR=build

settings:
	cp settings.example.json settings.dev.json

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

logs:
	docker compose -f docker-compose.yaml logs -f

clean:
	rmdir $(BUILD_DIR)

.PHONY: build run docker-build
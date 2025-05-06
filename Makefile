include .env
export
BINARY_NAME=myapp

## build: Build binary
build:
	@echo "Building..."
	env CGO_ENABLED=0  go build -ldflags="-s -w" -o ${BINARY_NAME} ./cmd/web
	@echo "Built!"

## run: go run
run:
	@echo "Starting..."
	go run ./cmd/web

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME} || true
	@rm tmp/*.pdf || true
	@echo "Cleaned!"

## start: build and run compiled app
start: build
	@echo "Starting..."
	@env ./${BINARY_NAME} &
	@echo "Started!"

## stop: stops the running application
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start

## test: runs all tests
test:
	go test -v ./...

new_migration:
	migrate create -ext sql -dir data/migration -seq ${name}

migrate_up:
	migrate -path data/migration -database ${DB_URL} -verbose up

migrate_up1:
	migrate -path data/migration -database ${DB_URL} -verbose up 1

migrate_down:
	migrate -path data/migration -database ${DB_URL} -verbose down

migrate_down1:
	migrate -path data/migration -database ${DB_URL} -verbose down 1

sqlc:
	sqlc generate

mock:
	mockgen -package mockdb -destination data/mock/store.go github.com/dubass83/go-concurrency-project/data/sqlc Store

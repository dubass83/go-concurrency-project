BINARY_NAME=myapp
DSN="host=localhost port=5432 user=postgres password=password dbname=concurrency sslmode=disable timezone=UTC connect_timeout=5"
DB_URL="postgresql://postgres:password@localhost:5432/concurrency?sslmode=disable"
REDIS="127.0.0.1:6379"

## build: Build binary
build:
	@echo "Building..."
	env CGO_ENABLED=0  go build -ldflags="-s -w" -o ${BINARY_NAME} ./cmd/web
	@echo "Built!"

## run: builds and runs the application
run: build
	@echo "Starting..."
	@env DSN=${DSN} REDIS_URL=${REDIS} ./${BINARY_NAME} &
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

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

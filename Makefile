## test: run test coverage on all package
test:
	@go test -v -cover -coverprofile=coverage.out -p 1 ./... 
	@go tool cover -func coverage.out | grep total

## docker-compose: build and run the docker compose
docker-compose: 
	echo "Pulling latest version"
	@docker-compose pull
	echo "Stop and remove the containers"
	@docker-compose rm -f
	echo "running the docker-compose container..."
	@docker-compose up -d --build

## mock: run mockgen on all package
mock:
	@go generate ./...
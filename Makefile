
# Which architecture to build
ARCH ?= amd64

# The docker image name
IMAGE := dealltest

REGISTRY := andryhardiyanto

# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)

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

## build docker image
docker: Dockerfile
	echo "building the $(IMAGE) container..."
	docker build --build-arg "VERSION=$(VERSION)" --label "version=$(VERSION)" -t $(IMAGE):$(VERSION) .

## push docker image to dockerhub registry
push-docker: .push-docker
.push-docker:
	docker tag $(IMAGE):$(VERSION) $(REGISTRY)/$(IMAGE):$(VERSION)
	docker push $(REGISTRY)/$(IMAGE):$(VERSION)
	echo "pushed: $(REGISTRY)/$(IMAGE):$(VERSION)"

## deploy k8s
deploy-k8s: 
	echo "apply user-deployment"
	kubectl apply -f scripts/deployment.yaml
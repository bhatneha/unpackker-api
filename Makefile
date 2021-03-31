GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')
PROJECT_NAME?=unpackker_api
APP_DIR?=$$(git rev-parse --show-toplevel)
VERSION?=0.0.1

.PHONY: help
help: ## Prints help (only for targets with comments)
	@grep -E '^[a-zA-Z0-9._-]+:.*?## .*$$' ${MAKEFILE_LIST} | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

local.fmt: ## Formats all the go code in th application
	gofmt -w ${GOFMT_FILES}

local.vend: local.fmt ## Loads all the dependencies required to vendor directory
	go mod vendor
	go mod tidy

local.build: local.vend ## loads the dependencies and Generates the artifact with the help of 'go build'
	go build -o ${PROJECT_NAME}

local.run: local.build ## Generates the artifact and Runs with the help of 'go run'
	./${PROJECT_NAME}

golang.lint: ## Lints application for errors, it is a linters aggregator (https://github.com/golangci/golangci-lint).
	docker run --rm -v ${APP_DIR}:/app -w /app golangci/golangci-lint:v1.27-alpine golangci-lint run --color always

dockerise: docker.lint golang.lint 1`## Containerise the appliction
	docker build -t ${DOCKER_USER}/${PROJECT_NAME}:${VERSION} .

docker.lint: ## Linting Dockerfile
	docker run --rm -v ${APP_DIR}:/app -w /app hadolint/hadolint:latest-alpine hadolint Dockerfile

docker.publish.image: docker.login ## Publisies the image to the registered docker registry.
	docker push ${DOCKER_USER}/${PROJECT_NAME}:${VERSION}

docker.login: ## Establishes the connection to the docker registry
	docker login -u ${DOCKER_USER} -p ${DOCKER_PASSWD} ${DOCKER_REPO}



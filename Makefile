GOCMD=go
MAIN_FILE?=cmd/app/main.go
BINARY_NAME?=server
APP_NAME?=petstore
COVEROUT?=cover.txt

.PHONY: build
build:
	CGO_ENABLED=0 ${GOCMD} build -v -o ${BINARY_NAME} ${MAIN_FILE}

.PHONY: run
run: build
	./${BINARY_NAME}

.PHONY: clean
clean:
	${GOCMD} clean
	rm ${BINARY_NAME}

.PHONY: test
test:
	${GOCMD} test -race ./... -coverpkg=./... -coverprofile=${COVEROUT}

${COVEROUT}: test

.PHONY: cover
cover: ${COVEROUT}
	${GOCMD} tool cover -html=${COVEROUT}

.PHONY: deps
deps: go.mod go.sum
	${GOCMD} mod download

.PHONY: generate
generate:
	${GOCMD} generate ./...

.PHONY: docker-build
docker-build: ${DOCKERFILE}
	docker build --tag ${APP_NAME} .
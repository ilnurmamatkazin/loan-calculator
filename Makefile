PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
$(shell [ -f bin ] || mkdir -p $(PROJECT_BIN))
PATH := $(PROJECT_BIN):$(PATH)

GOLANGCI_LINT = $(PROJECT_BIN)/golangci-lint

.PHONY: .install-linter
.install-linter:
	### INSTALL GOLANGCI-LINT ###
	[ -f $(PROJECT_BIN)/golangci-lint ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.58.1

.PHONY: lint
lint: .install-linter
	### RUN GOLANGCI-LINT ###
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

.PHONY: lint-fast
lint-fast: .install-linter
	$(GOLANGCI_LINT) run ./... --fast --config=./.golangci.yml

.PHONY: run
run:
	go run ./cmd/main.go -config=./cmd/config.yml

.PHONY: test-run
test-run:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile cover.out && go tool cover -func cover.out	

DOCKER_IMAGE = loan-calculator:local

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run:
	docker run $(DOCKER_IMAGE)

.PHONY: docker-stop
docker-stop:
	docker stop $$(docker ps -q --filter ancestor=$(DOCKER_IMAGE))

.PHONY: docker-rm
docker-rm:
	docker rm $$(docker ps -a -q --filter ancestor=$(DOCKER_IMAGE))
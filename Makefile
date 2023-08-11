
PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = v$(shell cat .version)
COMMIT_SHA ?= $(shell git describe --always)-devel

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBUILD=CGO_ENABLED=0 go build -a -ldflags "-X ${PKG}/version.Version=${VERSION}+sha.${COMMIT_SHA}"

WORKSPACE ?= webapp

up: tidy
	cd ./cmd/$(WORKSPACE) && go run .
clean:
	rm -rf ./cmd/$(WORKSPACE)/out


tidy:
	go mod tidy

upgrade:
	go get -u ./...

build:
	$(MAKE) build.webapp GOOS=linux GOARCH=amd64
	$(MAKE) build.webapp GOOS=linux GOARCH=arm64
	$(MAKE) build.webapp GOOS=darwin GOARCH=amd64
	$(MAKE) build.webapp GOOS=darwin GOARCH=arm64

build.webapp:
	@echo "Building webapp for $(GOOS)/$(GOARCH)"
	cd ./cmd/$(WORKSPACE) && $(GOBUILD) -o ../../out/webapp-$(GOOS)-$(GOARCH)

install: build.webapp
	mv ./out/webapp-$(GOOS)-$(GOARCH) ${GOPATH}/bin/webapp

release:
	git push
	git push origin ${VERSION}

webapp.config:
	COMMIT_SHA=${COMMIT_SHA} webapp config

webapp.buildx:
	COMMIT_SHA=${COMMIT_SHA} webapp buildx --with-builder --push

webapp.deploy:
	COMMIT_SHA=${COMMIT_SHA} webapp deploy

apply: install webapp.buildx webapp.deploy

debug: install
	COMMIT_SHA=${COMMIT_SHA} webapp run env | grep PROJECT_

docker:
	docker build -t example.com/demo/opentelemetry-gin-demo .

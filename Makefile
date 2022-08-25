# Paths to tools needed in dependencies
GO := $(shell which go)
DOCKER := $(shell which docker)

# nginx image and version
NGINX := "library/nginx"
VERSION := "1.23.1"
IMAGE := "nginx-gateway"

# target architectures: linux/amd64 linux/arm/v7 linux/arm64/v8

# Paths to locations, etc
BUILD_DIR := "build"
CMD_DIR := $(wildcard cmd/*)
BUILD_FLAGS := ""

# Targets
all: clean cmd

cmd: $(CMD_DIR)

test:
	@${GO} mod tidy
	@${GO} test -v ./pkg/...

$(CMD_DIR): dependencies mkdir FORCE
	@echo Build cmd $(notdir $@)
	@${GO} build ${BUILD_FLAGS} -o ${BUILD_DIR}/$(notdir $@) ./$@

docker: dependencies docker-dependencies
	@${DOCKER} build --tag ${IMAGE}-arm:${VERSION} --build-arg VERSION=${VERSION} --build-arg PLATFORM=linux/arm/v7 etc/docker
	@${DOCKER} build --tag ${IMAGE}-arm64:${VERSION} --build-arg VERSION=${VERSION} --build-arg PLATFORM=linux/arm64/v8 etc/docker
	@${DOCKER} build --tag ${IMAGE}-amd64:${VERSION} --build-arg VERSION=${VERSION} --build-arg PLATFORM=linux/amd64 etc/docker
	@${DOCKER} manifest create ${IMAGE}:${VERSION} --amend ${IMAGE}-arm:${VERSION} --amend ${IMAGE}-arm64:${VERSION} --amend ${IMAGE}-amd64:${VERSION}
	@${DOCKER} manifest annotate ${IMAGE}:${VERSION} ${IMAGE}-arm:${VERSION} --arch arm --os linux --variant v7
	@${DOCKER} manifest annotate ${IMAGE}:${VERSION} ${IMAGE}-arm64:${VERSION} --arch arm64 --os linux --variant v8
	@${DOCKER} manifest annotate ${IMAGE}:${VERSION} ${IMAGE}-amd64:${VERSION} --arch amd64 --os linux

FORCE:

docker-dependencies:
ifeq (,${DOCKER})
        $(error "Missing docker binary")
endif

dependencies:
ifeq (,${GO})
        $(error "Missing go binary")
endif

mkdir:
	@echo Mkdir ${BUILD_DIR}
	@install -d ${BUILD_DIR}

clean:
	@echo Clean
	@rm -fr $(BUILD_DIR)
	@${GO} mod tidy
	@${GO} clean


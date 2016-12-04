.PHONY: default all prepare build golint build-no-golint product integration-test clean

PWD := $(shell pwd)
BUNDLES := $(shell pwd)/bundles
JOB_NAME := $(if $(JOB_NAME),$(JOB_NAME),test)
BUILD_IMG := "mailer-build:latest"
PRODUCT_IMG := "mailer:latest"
DEPLOY_NAME := "mailer"
DEPLOY_DBDIR := "/var/lib/mailer/db"

default: build

all: integration-test push

prepare:
	mkdir -p $(BUNDLES)
	docker build --force-rm -t $(BUILD_IMG) -f Dockerfile.build .

golint: prepare
	docker run --rm -v $(PWD):/src:ro -e GOLINT_ONLY=yes $(BUILD_IMG)

build: prepare
	docker run --rm -v $(PWD):/src:ro -v $(BUNDLES):/product:rw  $(BUILD_IMG)

build-no-golint: prepare
	docker run --rm -v $(PWD):/src:ro -v $(BUNDLES):/product:rw -e NO_GOLINT=yes $(BUILD_IMG)

product: build
	mkdir -p $(BUNDLES)/latest/
	cp -avfL $(BUNDLES)/mailer-latest $(BUNDLES)/latest/mailer
	docker build --force-rm -t $(PRODUCT_IMG) -f Dockerfile.product .

deploy: product
	mkdir -p $(DEPLOY_DBDIR)
	docker rm -f $(DEPLOY_NAME) || true
	docker run --name=$(DEPLOY_NAME) --dns=114.114.114.114 -d -p 80:80 -p 127.0.0.1:27017:27017 -v $(DEPLOY_DBDIR):/mailer/db  $(PRODUCT_IMG)

integration-test: build
	echo "not implement yet"; exit 1
	docker run --rm --privileged -v $(TMP_DIR):/var/lib/docker:rw \
		-v $(BUNDLES):/bundles:rw csphere-build tools/integration-test.sh

clean:
	rm -rfv $(BUNDLES)


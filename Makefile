SHELL            := /bin/bash
MAKEFILE_DIR     := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
DIST_DIR         := $(abspath $(MAKEFILE_DIR)dist)
GO_DIST_DIR      := $(abspath $(MAKEFILE_DIR)go/dist)
KMTTG_VERSION    ?= v0.0.1
KMTTG_PORT       ?= 7676
MOUNT_DIR        = /mnt/kmttg
DOCKER_IMAGE_TAG ?= local
DOCKER_IMAGE     = tartale/kmttg-plus:$(DOCKER_IMAGE_TAG)
DOCKER_RUN_ARGS  = --rm --network=host --name kmttg-plus --platform linux/amd64 \
	-v /var/run/dbus:/var/run/dbus \
	-v $(CURDIR)/output:$(MOUNT_DIR)/output:rw \
	-e KMTTG_PORT=$(KMTTG_PORT) -e KMTTG_MEDIA_ACCESS_KEY

NVM_DIR ?= $(HOME)/.nvm

clean:
	docker rmi $(DOCKER_IMAGE) || true

build:
	cd go; \
	make build

test:
	cd go; \
	make test

install:
	cd go && $(MAKE) install DIST_DIR=$(DIST_DIR)

image: webui
	cd go && $(MAKE) install GOARCH=amd64 GOOS=linux DIST_DIR=$(DIST_DIR)
	docker build --build-arg KMTTG_VERSION=$(KMTTG_VERSION) --platform linux/amd64 -t $(DOCKER_IMAGE) .

image-run:
	docker run $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE)

image-bg:
	docker run -d $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE)

push:
	docker push $(DOCKER_IMAGE)

run:
	cd go; \
	make run

shell:
	docker run -it $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE) /bin/bash

watch:
	trap 'kill 0' EXIT; \
	(cd go && $(MAKE) dev) & \
	(cd webui && \
		{ [ ! -s "$(NVM_DIR)/nvm.sh" ] || { source "$(NVM_DIR)/nvm.sh" && nvm install; }; } && \
		npm install && \
		BROWSER=none npm start) & \
	wait

webui:
	cd webui && \
	{ [ ! -s "$(NVM_DIR)/nvm.sh" ] || { source "$(NVM_DIR)/nvm.sh" && nvm install; }; } && \
	npm install && \
	npm run build && \
	DIST_DIR=$(GO_DIST_DIR) npm run deploy

.PHONY: clean build test install image image-run image-bg push run shell watch webui

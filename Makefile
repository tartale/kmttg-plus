DOCKER_IMAGE_TAG ?= local
DOCKER_IMAGE = tartale/kmttg-plus:$(DOCKER_IMAGE_TAG)
DOCKER_RUN_ARGS = --rm -v $(CURDIR)/overrides:$(MOUNT_DIR)/overrides -v $(CURDIR)/output:$(MOUNT_DIR)/output:rw -p 8181:8181

MOUNT_DIR = /mnt/kmttg

all: java image push run shell

java:
	cd java; \
	ant release

go:
	cd go; \
	make build

image:
	docker build -t $(DOCKER_IMAGE) .

push: image
	docker push $(DOCKER_IMAGE)

run:
	docker run -d $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE)

shell:
	docker run -it $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE) /bin/bash

.PHONY: all java go image push run shell

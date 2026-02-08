DOCKER_IMAGE_TAG ?= local
DOCKER_IMAGE = tartale/kmttg-plus:$(DOCKER_IMAGE_TAG)
DOCKER_RUN_ARGS = -it --rm -v $(CURDIR)/overrides:$(MOUNT_DIR)/overrides -v $(CURDIR)/output:$(MOUNT_DIR)/output:rw -p 8181:8181
DOCKER_RUN_CMD = docker run $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE)

MOUNT_DIR = /mnt/kmttg

all: java image push run shell

java:
	cd java; \
	ant release

image:
	docker build -t $(DOCKER_IMAGE) .

push: image
	docker push $(DOCKER_IMAGE)

run:
	$(DOCKER_RUN_CMD)

shell:
	$(DOCKER_RUN_CMD) /bin/bash

.PHONY: all java image push run shell

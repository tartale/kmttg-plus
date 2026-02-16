KMTTG_VERSION ?= v2.9.4-l
MOUNT_DIR = /mnt/kmttg
DOCKER_IMAGE_TAG ?= local
DOCKER_IMAGE = tartale/kmttg-plus:$(DOCKER_IMAGE_TAG)
DOCKER_RUN_ARGS = --rm -v $(CURDIR)/overrides:$(MOUNT_DIR)/overrides -v $(CURDIR)/output:$(MOUNT_DIR)/output:rw -p 8181:8181

clean:
	docker rmi $(DOCKER_IMAGE)
	cd java; \
	ant clean

java-build:
	cd java; \
	ant release

java-run:
	cd java; \
	./release/kmttg

go-build:
	cd go; \
	make build

go-run:
	cd go; \
	make run

image:
	docker build --build-arg KMTTG_VERSION=$(KMTTG_VERSION) -t $(DOCKER_IMAGE) .

image-run:
	docker run -d $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE)

push: image
	docker push $(DOCKER_IMAGE)

shell:
	docker run -it $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE) /bin/bash

.PHONY: clean java-build java-run go-build go-run image image-run push shell

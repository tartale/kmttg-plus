MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
KMTTG_VERSION ?= v0.0.1
MOUNT_DIR = /mnt/kmttg
DOCKER_IMAGE_TAG ?= local
DOCKER_IMAGE = tartale/kmttg-plus:$(DOCKER_IMAGE_TAG)
DOCKER_RUN_ARGS = --rm --name kmttg-plus -h kmttg.local -p $(KMTTG_PORT):$(KMTTG_PORT) \
	-v $(CURDIR)/overrides:$(MOUNT_DIR)/overrides -v $(CURDIR)/output:$(MOUNT_DIR)/output:rw \
	-e KMTTG_PORT -e KMTTG_MEDIA_ACCESS_KEY 

# 	-v /var/run/dbus:/var/run/dbus -v /var/run/avahi-daemon/socket:/var/run/avahi-daemon/socket \ /sys/fs/cgroup/openrc.dbu

clean:
	docker rmi $(DOCKER_IMAGE) || true
	rm -rf "$(KMTTG_CACHE_DIR)/*"
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

go-test:
	cd go; \
	go test ./...

go-run:
	cd go; \
	make run

go-install:
	cd go; \
	GOARCH=amd64 GOOS=linux GOBIN=$(MAKEFILE_DIR)/dist go install ./cmd/kmttg.go

image:
	docker build --build-arg KMTTG_VERSION=$(KMTTG_VERSION) -t $(DOCKER_IMAGE) .

image-run:
	docker run $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE)

image-bg:
	docker run -d $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE)

push: image
	docker push $(DOCKER_IMAGE)

shell:
	docker run -it $(DOCKER_RUN_ARGS) $(DOCKER_IMAGE) /bin/bash

.PHONY: clean java-build java-run go-build go-run image image-run push shell

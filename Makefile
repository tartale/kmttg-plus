.PHONY: check-env build push run-local run-image debug-image

check-env:
ifeq ($(APP_DIR),)
	$(error "APP_DIR env variable must be set to the root directory of the kmttg app")
endif
ifeq ($(MOUNT_DIR),)
	$(error "MOUNT_DIR env variable must be set to the directory to be mounted into the docker container")
endif
ifeq ($(TOOLS_DIR),)
	$(error "TOOLS_DIR env variable must be set to the root directory of the various tools (comskip, etc) that kmttg uses")
endif

MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

build:
	docker build -t tartale/kmttg-plus:latest .

push: build
	docker push tartale/kmttg-plus:latest

run-local: check-env
	$(MAKEFILE_DIR)/kmttg.sh

run-image: check-env build
	echo "running kmttg container on the local machine"; \
  docker run --network=host --rm \
		-v $(MOUNT_DIR):/mnt/kmttg \
		--platform=linux/amd64 tartale/kmttg-plus:latest

debug-image: check-env build
	echo "debugging kmttg container on the local machine"; \
  docker run -it --network=host --rm \
		-v $(MOUNT_DIR):/mnt/kmttg \
		--platform=linux/amd64 tartale/kmttg-plus:latest /bin/bash

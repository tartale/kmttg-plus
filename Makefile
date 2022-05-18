.PHONY: check-env build push run-local run-image debug-image

check-env:
ifeq ($(APP_DIR),)
	$(error "APP_DIR env variable must be set to the root directory of the kmttg app")
endif
ifeq ($(OUTPUT_DIR),)
	$(error "OUTPUT_DIR env variable must be set to the root directory for all output files")
endif
ifeq ($(TOOLS_DIR),)
	$(error "TOOLS_DIR env variable must be set to the root directory of the various tools (comskip, etc) that kmttg uses")
endif

MAKEFILE_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

build:
	docker build -t tartale/kmttg-plus:latest .

push: build
	docker push tartale/kmttg-plus:latest

output-paths: check-env
	echo "creating output files/directories in $(OUTPUT_DIR)"; \
	mkdir -p $(OUTPUT_DIR)/mpeg && \
	mkdir -p $(OUTPUT_DIR)/encode && \
	mkdir -p $(OUTPUT_DIR)/files && \
	mkdir -p $(OUTPUT_DIR)/qsfix && \
	mkdir -p $(OUTPUT_DIR)/webcache && \
	mkdir -p $(OUTPUT_DIR)/mpegcut && \
	mkdir -p $(OUTPUT_DIR)/output && \
	touch $(OUTPUT_DIR)/auto.history
	touch $(OUTPUT_DIR)/auto.log.0

run-local: check-env output-paths
	echo "running kmttg app on the local machine"; \
	envsubst < $(MAKEFILE_DIR)/config.ini.personal > $(APP_DIR)/config.ini && \
	envsubst < $(MAKEFILE_DIR)/config.ini.template >> $(APP_DIR)/config.ini && \
	cp -f $(MAKEFILE_DIR)/auto.ini $(APP_DIR)/auto.ini && \
	ln -f -s $(OUTPUT_DIR)/auto.history $(APP_DIR)/auto.history && \
	ln -f -s $(OUTPUT_DIR)/auto.log.0 $(APP_DIR)/auto.log.0 && \
	cd $(APP_DIR); ./kmttg

run-image: check-env build output-paths
	echo "running kmttg container on the local machine"; \
  docker run --network=host --rm -v $(OUTPUT_DIR):/mnt/kmttg/output --platform=linux/amd64 tartale/kmttg-plus:latest

debug-image: check-env build output-paths
	echo "debugging kmttg container on the local machine"; \
  docker run -it --network=host --rm -v $(OUTPUT_DIR):/mnt/kmttg/output --platform=linux/amd64 tartale/kmttg-plus:latest /bin/bash

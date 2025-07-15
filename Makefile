# Variables (default values, podem ser sobrescritas na linha de comando)
IMAGE_NAME ?= myapp
VERSION ?= 1.0.0

# Default target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build           Run build.sh script (uses IMAGE_NAME and VERSION)"
	@echo "  deploy          Run deploy.sh script (uses IMAGE_NAME and VERSION)"
	@echo "  mocks           Run generate-mocks.sh script (pass args as MOCK_ARGS)"
	@echo "  test            Run test.sh script"
	@echo "  setup           Run setup.sh script"
	@echo "  swagger         Run swagger-doc.sh script"

.PHONY: build
build:
	@scripts/build.sh $(IMAGE_NAME) $(VERSION)

.PHONY: deploy
deploy:
	@scripts/deploy.sh $(IMAGE_NAME) $(VERSION)

.PHONY: mocks
mocks:
	@scripts/generate-mocks.sh $(MOCK_ARGS)

.PHONY: test
test:
	@scripts/test.sh

.PHONY: setup
setup:
	@scripts/setup.sh

.PHONY: swagger
swagger:
	@scripts/swagger-doc.sh

#!/bin/bash

set -e

echo "Starting setup..."

# Check Go installation
if ! command -v go >/dev/null 2>&1; then
  echo "Go is not installed. Please install Go and rerun this script."
  exit 1
fi

# Check Docker installation
if ! command -v docker >/dev/null 2>&1; then
  echo "Docker is not installed. Please install Docker and rerun this script."
  exit 1
fi

# Install Uber mockgen if not installed
if ! command -v mockgen >/dev/null 2>&1; then
  echo "Installing Uber's mockgen..."
  go install github.com/uber/mock/mockgen@latest
fi

# Optionally install golangci-lint
if ! command -v golangci-lint >/dev/null 2>&1; then
  echo "Installing golangci-lint..."
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.3
fi

# Create folders
echo "Creating folders 'mocks' and 'bin' if not exist..."
mkdir -p mocks bin

echo "Setup finished."

echo ""
echo "Add GOPATH/bin to your PATH if it's not there already:"
echo "  export PATH=\$PATH:$(go env GOPATH)/bin"
echo ""
echo "You can now run scripts like generate-all-mocks.sh and others."

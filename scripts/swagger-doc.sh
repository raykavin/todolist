#!/usr/bin/env bash
set -e

echo "Instalando swag CLI..."
go install github.com/swaggo/swag/cmd/swag@latest

echo "Gerando documentação Swagger..."
swag init --parseDependency --parseInternal -g cmd/api/main.go

echo "Documentação gerada na pasta ./docs"

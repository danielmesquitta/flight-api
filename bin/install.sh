#!/bin/bash

packages=(
    "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2"
    "github.com/segmentio/golines@latest"
    "github.com/air-verse/air@latest"
    "github.com/google/wire/cmd/wire@latest"
    "github.com/danielmesquitta/wire-config@latest"
    "github.com/swaggo/swag/cmd/swag@latest"
    "github.com/zxmfke/swagger2openapi3/cmd/swag2op@latest"
    "github.com/vektra/mockery/v3@v3.2.2"
)

echo "Installing and updating Go packages..."

for package in "${packages[@]}"; do
    echo "$package..."
    go install "$package"
done

echo "All packages have been successfully installed and updated."

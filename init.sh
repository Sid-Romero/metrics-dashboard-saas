#!/bin/bash

# Core files
touch README.md docker-compose.yml

# Backend
mkdir -p backend/handlers
touch backend/main.go backend/go.mod backend/handlers/metrics.go

# Agent
mkdir -p agent/system
touch agent/main.go agent/go.mod agent/system/collector.go

# Infra
mkdir -p infra/{terraform,ansible}

# Monitoring
mkdir -p monitoring/{grafana,prometheus}

# Scripts
mkdir -p scripts
touch scripts/{build.sh,deploy.sh}

# GitHub Actions
mkdir -p .github/workflows
touch .github/workflows/ci.yml

echo "✅ Project structure for '$PROJECT' created."

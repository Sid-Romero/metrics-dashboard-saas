#!/bin/bash

# Dossier chart
mkdir -p chart/metrics-saas/templates

# Chart.yaml
cat <<EOF > chart/metrics-saas/Chart.yaml
apiVersion: v2
name: metrics-saas
description: A metrics dashboard SaaS (DB + Backend + Agent + Grafana)
type: application
version: 0.1.0
appVersion: "1.0"
EOF

# values.yaml
cat <<EOF > chart/metrics-saas/values.yaml
namespace: metrics

postgres:
  user: metrics
  password: metrics
  db: metricsdb
  storage: 1Gi

backend:
  image: sidya18/metrics-backend:latest
  port: 8080

agent:
  image: sidya18/metrics-agent:latest

grafana:
  adminUser: admin
  adminPassword: supersecret
  port: 3000

retention: 12 hours
EOF

# Templates vides
for f in secrets db-statefulset db-service backend-deployment backend-service agent-daemonset grafana-deployment grafana-service grafana-datasource-cm grafana-dashboard-cm cronjob-cleaner; do
    touch chart/metrics-saas/templates/$f.yaml
done

echo "✅ Helm chart skeleton created under chart/metrics-saas/"

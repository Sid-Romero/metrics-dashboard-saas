# Metrics Dashboard SaaS

This project is a cloud-native monitoring platform that collects and visualizes system metrics (CPU, memory, etc.) from distributed agents.

## Features
- Lightweight agent written in Go that collects CPU and memory usage.
- API server (Go) that receives and stores metrics.
- PostgreSQL/TimescaleDB backend for time-series data.
- Grafana dashboards for visualization.
- Infrastructure as Code with Terraform (AWS EC2, RDS).
- Configuration management with Ansible for agent deployment.
- CI/CD pipelines with GitHub Actions.
- Containerized with Docker, deployable on Kubernetes.

## Architecture
1. Agents collect metrics and send them to the backend API.
2. API stores metrics into PostgreSQL.
3. Grafana connects to the database for visualization.
4. Terraform provisions cloud infrastructure.
5. Ansible automates deployment of agents on servers.

## Tech Stack
- **Language:** Go
- **Database:** PostgreSQL / TimescaleDB
- **Visualization:** Grafana
- **Infra:** Terraform, Ansible, AWS
- **Orchestration:** Docker, Kubernetes
- **CI/CD:** GitHub Actions

## Getting Started
1. Clone the repository
   ```bash
   git clone https://github.com/<your-username>/metrics-dashboard-saas.git
   cd metrics-dashboard-saas
   ```

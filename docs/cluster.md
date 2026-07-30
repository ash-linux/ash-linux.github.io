<![CDATA[# Cluster Mode

Scale Ash across multiple machines for distributed AI inference.

> **This is an advanced feature.** Most users only need a single Ash instance. Cluster mode is for users who want to distribute AI workloads across multiple nodes.

## What It Does

The Ash cluster controller spins up multiple Ollama instances behind a load balancer. This lets you:

- **Run larger models** by splitting across nodes
- **Handle more concurrent requests** with load balancing
- **Distribute embedding workloads** across multiple machines

## Quick Start

```bash
# Spin up a 3-node cluster
ash-cluster up --nodes 3

# Check status
ash-cluster status

# View all running containers
ash-cluster ps

# View logs
ash-cluster logs

# Tear down
ash-cluster down
```

## Architecture

```
                    ┌──────────────┐
                    │   Traefik    │  ← Load balancer + dashboard (:8080)
                    │   Router     │
                    └──────┬───────┘
                           │
           ┌───────────────┼───────────────┐
           │               │               │
    ┌──────▼──────┐ ┌──────▼──────┐ ┌──────▼──────┐
    │  Ollama     │ │  Ollama     │ │  Ollama     │
    │  Node 1     │ │  Node 2     │ │  Node 3     │
    └─────────────┘ └─────────────┘ └─────────────┘
           │
    ┌──────▼──────┐        ┌──────────────┐
    │   Qdrant    │        │   Consul     │  ← Service discovery (:8500)
    │   Vector DB │        │              │
    └─────────────┘        └──────────────┘
```

## Components

| Service | Port | Role |
|---------|------|------|
| **Traefik** | `:8080` | Load balancer dashboard |
| **Router** (Nginx) | `:80`, `:11434` | Reverse proxy, routes requests |
| **Consul** | `:8500` | Service discovery and health checking |
| **Qdrant** | `:6333` | Shared vector database |
| **Ollama Nodes** | Internal | AI model inference workers |

## Orchestrators

### Docker Compose (Default)

```bash
ash-cluster up --nodes 3
```

Uses Docker Compose to manage all services. Best for local development and single-machine clusters.

### HashiCorp Nomad

```bash
ash-cluster up --nodes 5 --orchestrator nomad
```

Uses Nomad for job scheduling. Better for multi-machine production deployments.

## Using the Cluster

Once running, use the same Ollama API on `localhost:11434` — the router automatically load-balances requests:

```bash
# Pull a model on a specific node
docker exec ash-ollama-node-1 ollama pull llama3.2

# Query through the load balancer (hits any available node)
curl http://localhost:11434/api/generate \
  -d '{"model":"llama3.2","prompt":"Hello"}'
```

## Cloud Deployment

Deploy Ash clusters on cloud infrastructure using the included Terraform configs:

```bash
# AWS
cd terraform/aws-ash-cluster && terraform init && terraform apply

# GCP
cd terraform/gcp-ash-cluster && terraform init && terraform apply

# Azure
cd terraform/azure-ash-cluster && terraform init && terraform apply
```

## Requirements

- Docker and Docker Compose installed
- At least 4 GB RAM per node
- Consul (optional, for service discovery)
- Nomad (optional, for Nomad orchestrator)

---

**Next:** [AI Agent Integration →](ai-agents.md) | [GPU & Performance →](gpu-and-performance.md)
]]>

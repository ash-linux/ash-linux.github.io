<![CDATA[# Changelog

All notable changes to Ash Linux are documented here.

## v2.0.1 (2026-07-20)

### Documentation
- Complete rewrite of all documentation for end users
- New guides: How Search Works, VMware Setup, Backup & Restore, Troubleshooting, Health Checks, Cluster Mode, AI Agent Integration
- Updated README with user-focused feature descriptions, FAQ, and platform matrix
- Removed all developer-only and aspirational feature references

---

## v2.0.0 (2026-07-20)

### Features
- **Semantic file search** — Press `Super+Space` to search files by meaning using vector embeddings
- **Pure-bash launcher hook** — ~50 lines of bash, zero Python dependencies in the search path
- **nomic-embed-text** — 768-dim Ollama-native embedding model, fast on CPU
- **Qdrant standalone binary** — No Docker, no AUR, no Python SDK. Downloaded directly from GitHub
- **LSFS indexing daemon** — Python background service that watches files and auto-indexes into Qdrant
- **Auto-login + auto-start** — Boots straight into Hyprland with all services running
- **Time-based fallback** — Queries like "files from 2h" use `fd`/`find` when vector search isn't needed
- **One-shot installer** — `ultimate-fix-v2.sh`, idempotent and safe to re-run

### Infrastructure
- VMware display and clipboard workarounds documented
- systemd sandboxing for Qdrant and LSFS daemon
- Multi-node cluster controller (`ash-cluster`)
- Packer templates for VMware, AWS, GCP, Azure, Vagrant
- Terraform configs for cloud cluster deployment
- Multi-channel distribution (Cloudflare R2, Bunny CDN, Archive.org, GHCR, Docker Hub, Quay.io)
- AI agent configurator for Claude, Cursor, Windsurf, Cline, Gemini
- GitHub Actions CI/CD pipelines

---

## v1.x (2026-07-19)

### Initial Release
- Initial deployment scripts
- Basic Arch Linux ISO profile
]]>

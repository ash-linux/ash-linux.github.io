<![CDATA[# Changelog

All notable changes to Ash Linux are documented here.

## v4.1.0 (2026-07-30)

### Features
- **Remote index federation** — Query across peer machines using `ash-federate` and `ash_federation.py`.
- **Multi-node cluster controller** — Use `ash-cluster` to deploy Qdrant clusters.
- **Content extraction for binary files** — Semantic search within PDFs, DOCX, Images, and Audio via `ash_extract.py`.
- **File relationship graph** — Understand imports, symlinks, and dependencies via `ash-graph`.
- **Storage tiering** — Hot (memory), Warm (disk), Cold (archive) support configured in the installer.
- **Query handler plugin system** — Extend the desktop search with JSON plugins (e.g., GitHub, Calc) via `ash-plugin`.
- **OpenAPI support** — Documentation and `ash-api` HTTP server for standardized integration.
- **Webhook notifications** — Dispatch events to Slack, Discord, or ntfy from system events.
- **VSCode Extension** — Scaffolded `ash-search` to query Ash OS from within your editor.
- **agy Enhancements** — Added `agy health`, `agy log`, and `agy phase` commands.

---

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

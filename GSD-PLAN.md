# GSD-PLAN — Ash Linux Grand Strategy Document

> A debatable roadmap with 8 phases, each containing a key decision fork with trade-offs, recommendations, and reasoning. This is the document that defines the next 6-12 months.

**Status:** Active · **Owner:** os · **Last updated:** 2026-07-30

---

## Top 5 Debates That Define the Next 6 Months

### Debate 1: Native vs Container-First Install
**Recommended: Native with `--container` opt-in**

| Path | Approach | Trade-off |
|------|----------|-----------|
| **A — Native bare-metal** | Install packages, systemd units, binaries on the host | Full performance, hardware access, simple debugging; harder to clean up, distro-specific |
| **B — Container-first** | Everything in podman/Docker quadlets | Clean isolation, distro-agnostic, easy uninstall; perf overhead, GPU/compositor pain |
| **C — Hybrid** | Native by default, `--container` flag for Qdrant/Ollama | Best of both; two code paths to maintain |

**Decision:** Native default. `--container` flag lands in Phase 3. Rationale: this is a desktop OS — users expect native performance and hardware access. Container mode serves CI/server/SteamOS. The two paths share 70% of the install logic (config, health checks, backup).

### Debate 2: Hybrid Search Algorithm
**Recommended: Adaptive RRF with regex heuristics**

| Path | Approach | Trade-off |
|------|----------|-----------|
| **A — Pure semantic** | Query → embed → Qdrant search → results | Misses exact-match queries, bad for code/filesystem queries |
| **B — Keyword + semantic fused** | Run both searches, fuse with RRF | More latency, complex tuning |
| **C — Adaptive classifier** | Classify query type first, route to appropriate engine | Best accuracy; classifier is another model to maintain |

**Decision:** Adaptive RRF. Classify query with regex patterns (cheap, no model load), route to semantic or keyword or both — fuse results with reciprocal rank fusion. No new model dependency.

### Debate 3: Systemd Hardening Level
**Recommended: Moderate (don't break UX)**

| Path | Approach | Trade-off |
|------|----------|-----------|
| **A — Minimal** | Default systemd units, no sandboxing | Works everywhere; insecure |
| **B — Moderate** | `ProtectSystem=strict`, `NoNewPrivileges=true`, `PrivateTmp=true` | Good security; some services need path exceptions |
| **C — Maximum** | Full sandbox + SELinux/AppArmor + seccomp | Very secure; breaks until profiles are tuned per-distro |

**Decision:** Moderate (current state). Add `ProtectHome=read-only` for Qdrant. AppArmor profiles ship in Phase 5 but are opt-in until v5.0.

### Debate 4: Config Format
**Recommended: Declarative YAML with Ansible converter**

| Path | Approach | Trade-off |
|------|----------|-----------|
| **A — Shell env vars** | `/etc/ash/*.env`, simple key=value | Works now; no structure, no validation |
| **B — TOML** | Rust-ecosystem standard, strict schema | Good for tools; unfamiliar to non-Rust users |
| **C — YAML** | Universal, Ansible-compatible, schema-validatable | YAML has footguns; but enables `--export-ansible` |

**Decision:** YAML. Ship a schema (`/etc/ash/schema.yaml`), validate with `yq` or Python. `ash-install --export-ansible` generates a full Ansible role from the config. Endgame: `ash config set search.model=nomic-embed-text` CLI.

### Debate 5: Architecture — Language & Component Boundary
**Recommended: Bash orchestration + Go agents (strangler fig)**

| Path | Approach | Trade-off |
|------|----------|-----------|
| **A — All bash** | Everything in shell scripts | 1200-line installer works today; becomes unmaintainable at 5000+ lines |
| **B — All Python** | Python CLI tools, daemon already in Python | High memory overhead for CLI; dependency hell |
| **C — Strangler fig** | Keep bash for install/bootstrap, rewrite critical paths in Go | Best long-term; dual-language overhead during transition |

**Decision:** Strangler fig. The installer stays in bash (it's the bootstrap — it needs zero dependencies). New components (ash-doctor, ash-query, ash-config, agy) are written in Go. The LSFS daemon stays in Python (it's already working). Incrementally replace Python CLIs with Go as we add features.

---

## Phase 1 — Install & Bootstrap (Weeks 1-3)
**Theme:** Make the install bulletproof, cross-distro, and self-healing.

### Key Question: How much runtime state should the installer manage?

| Aspect | Path A — Stateless | Path B — Stateful | Path C — Hybrid |
|--------|-------------------|------------------|-----------------|
| State file | None — detect-and-run | Full JSON state machine | JSON for install, detect for repair |
| Re-run | Always full install | Resumes from last phase | Repair mode detects broken services |
| Features | Simple, 800 lines | Complex, 1800+ lines | Complex, 1500 lines |
| **recommended** | | | **→ C** |

### Deliverables
- [x] `--dry-run` flag
- [x] Preflight compatibility check (kernel, disk, RAM, ports)
- [x] Btrfs snapshot before install
- [x] State machine with resume from failure
- [x] Self-update with changelog diff
- [x] Cross-distro package translation (APT, DNF, Zypper, APK) — **Phase 1 priority**
- [x] `pm_translate` expansion — all package names mapped for 5 distro families
- [x] Offline bundle generation (`ash-install --bundle`)
- [x] Install telemetry (opt-in, anonymous)

### Decision Fork: Package Management
**Chosen:** Runtime detection with per-distro map files in `/etc/ash/pkgmaps/`.

### Risk
Distro-specific package names will drift. Mitigation: community-maintained maps via GitHub PRs, validated in CI.

---

## Phase 2 — Core Search Architecture (Weeks 4-6)
**Theme:** Make the search fast, accurate, and resource-conscious.

### Key Question: Single embedding model vs multi-model routing?

| Aspect | Path A — One Model | Path B — Model Switcher | Path C — Adaptive RRF |
|--------|-------------------|------------------------|----------------------|
| Complexity | Low | Medium | Medium-High |
| Accuracy | Good for semantic | Better for varied queries | Best for mixed workloads |
| Memory | ~140MB | ~500MB (2 models) | ~140MB + regex |
| **recommended** | | | **→ C** |

### Deliverables
- [x] nomic-embed-text integration
- [x] BM25 keyword fallback (SQLite FTS5 or tantivy)
- [x] Adaptive query classifier (regex-based, no model)
- [x] RRF fusion with tunable weights (`--search-mode hybrid|semantic|keyword`)
- [x] Inotify-based file watcher (replace 60s poll)
- [x] Compiled `.lsfsignore` (regex trie, not linear scan)
- [x] Batch embedding (16 chunks at once via `/api/embed`)
- [x] `lru_cache` on embedding calls to deduplicate
- [x] Index deduplication (SHA256 content hash → one embedding per unique file)

### Decision Fork: File Watcher Implementation
**Chosen:** `inotifywait` (inotify-tools) for initial scan + `watchdog` Python lib for persistent watching. No new system dependencies beyond what's already installed.

---

## Phase 3 — Vibecoding & Developer Experience (Weeks 7-9)
**Theme:** Make ASH the OS you want to code in, then throw away.

### Key Question: How disposable should the OS be?

| Aspect | Path A — Traditional | Path B — Ephemeral | Path C — Snapshotted |
|--------|---------------------|-------------------|---------------------|
| Workspace model | Files on disk, manual cleanup | Everything in `~/.ash/workspaces/` as btrfs snapshots | Hybrid — snapshots for projects, persistent for config |
| Reset | Reinstall OS | `ash reset --workspace` clears to snapshot | `ash workspace reset` per-project |
| **recommended** | | | **→ C** |

### Deliverables
- [x] `ash workspace create` — create a disposable dev environment
- [x] `ash workspace reset` — roll back to clean state
- [x] `ash up` / `ash down` — start/stop all services (like docker compose)
- [x] `ash-ask` — RAG over personal files (FUTURE #30)
- [x] Plugin system for query handlers (FUTURE #69)
- [x] Semantic clipboard history (FUTURE #29)
- [x] `Super+Space` sound effect (FUTURE #37)
- [x] Welcome wizard improvements — model selection, telemetry opt-in
- [x] TUI dashboard for workspace management

### Vibecoding Features
- [x] `agi` — "ash go install" — AI-suggested package installs
- [x] `agi` — "ash generate init" — scaffold a new project from natural language
- [x] Automatic `.gitignore` → `.lsfsignore` sync
- Git-aware search: "files I changed last Tuesday that touch authentication"

---

## Phase 4 — Container & Cross-Platform (Weeks 10-12)
**Theme:** Run ASH anywhere — Docker, non-Arch Linux, macOS, WSL2.

### Decision Fork: Container Strategy
**Chosen:** Rootless podman quadlets for Qdrant + Ollama. LSFS daemon stays native (needs filesystem access).

### Cross-Platform Support Matrix

| Platform | Phase 4 target | Status |
|----------|---------------|--------|
| Arch Linux | Current — stable | ✅ Now |
| Debian/Ubuntu | APT package mapping | 🔧 Phase 4 |
| Fedora | DNF package mapping | 🔧 Phase 4 |
| openSUSE | Zypper package mapping | 🔧 Phase 4 |
| Alpine | APK package mapping | 🔧 Phase 4 |
| Docker (rootless) | podman quadlet | 🔧 Phase 4 |
| macOS | Homebrew + launchd | 📅 Phase 5 |
| WSL2 | Windows interop | 📅 Phase 6 |

### Deliverables
- [x] `--container` install mode (podman quadlets)
- [x] Complete `pm_translate` for all 5 package managers
- [x] Docker/podman-compose stacks for server profile
- [x] CI/CD pipeline — build + test on all 5 distros
- [x] `ash-install --export-ansible` — generate Ansible role from config
- [x] Vagrant box provisioning support
- [x] Ephemeral CI mode (`--ephemeral` — install, test, uninstall, return exit code)

### Risk
Five distros will find five different bugs. Mitigation: CI runs on all five. Static analysis before merge.

---

## Phase 5 — Security & Hardening (Weeks 13-15)
**Theme:** Make ASH the most secure AI desktop OS, without breaking UX.

### Decision Fork: Security Level
**Chosen:** Moderate with opt-in maximum. The desktop user gets moderate (sandboxed systemd services, API key auth, firewall). The server/CI user can opt into maximum (AppArmor, auditd, TPM measurements).

### Deliverables
- [x] AppArmor profiles for qdrant, ollama, lsfs-daemon (FUTURE #10)
- [x] Secrets encryption via `age` + TPM key derivation (FUTURE #61)
- [x] GPG signing of install manifest (FUTURE #60)
- [x] Verified boot measurement (FUTURE #59)
- [x] Audit log (FUTURE #58) — append-only, `chattr +a`
- [x] Rate limiting on search endpoints
- [x] `ash-doctor --security-scan` — full security audit
- [x] USB/offline signing for air-gapped deploys (FUTURE #62)
- [x] SELinux policy module (optional, separate package)

### Security Model

```
User → Super+Space/CLI
  ↓                    ┌─────────────────┐
  └→ lsfs-query ──────→│  Audit Log      │
       ↓               │  (append-only)   │
  ┌────┴────┐          └─────────────────┘
  │Regex    │──miss──→│ Keyword FTS5    │
  │Classifier         │ (local only)    │
  └────┬────┘          └─────────────────┘
       │ hit
  ┌────┴──────────────────────────────┐
  │  Qdrant (localhost:6333)          │
  │  API key auth                     │
  │  Firewall: localhost only         │
  │  AppArmor: confined               │
  │  systemd: ProtectSystem=strict    │
  └───────────────────────────────────┘
```

---

## Phase 6 — Observability & TUI (Weeks 16-18)
**Theme:** Every component visible, every metric measurable, every problem fixable from a terminal.

### Decision Fork: TUI Framework
**Chosen:** `ash-tui` — a Go TUI app using `bubbletea` for real-time dashboards. Falls back to whiptail/dialog for installer. The TUI is not a web dashboard — it's a terminal app that feels native to the ASH experience.

### Deliverables
- [x] `ash-tui` — real-time dashboard (services, index status, resource usage)
- [x] `ash-doctor` — 40+ health checks, export to JSON (FUTURE #32)
- [x] `ash-stats` — storage usage dashboard (FUTURE #67)
- [x] HTML install report with waterfall chart (FUTURE #35)
- [x] Desktop notifications during long operations (FUTURE #41)
- [x] Live reload of configuration via SIGHUP (FUTURE #70)
- [x] Log rotation + archiving (FUTURE #36)
- [x] Split-pane TUI dashboard for install (FUTURE #31)

### Metrics to Surface
```
Services:  ● Qdrant :6333  ● Ollama :11434  ● LSFS daemon
Storage:   2.3 GB indexed (12,847 files)   Cache: 340 MB
Model:     nomic-embed-text (768-dim)   GPU: NVIDIA
Snapshots: 3 (using 2.1 GB)   Last backup: 3h ago
```

---

## Phase 7 — Federation & Scale (Weeks 19-21)
**Theme:** One ASH is useful. Many ASH machines are a distributed AI mesh.

### Decision Fork: Federation Protocol
**Chosen:** gRPC proxy with Qdrant's distributed deployment mode. Lightweight — no raft consensus for the index. Machines discover each other via mDNS or a config file.

### Deliverables
- [x] Remote index federation (FUTURE #68)
- [x] Multi-node cluster controller (existing `ash-cluster/`)
- [x] Hybrid local + remote model support (FUTURE #27)
- [x] File relationship graph (FUTURE #28)
- [x] Content extraction for binary files — PDF, DOCX, images, audio (FUTURE #26)
- [x] Storage tiering — hot/warm/cold (FUTURE #63)
- [x] Btrfs auto-defrag for Qdrant storage (FUTURE #66)

### Architecture

```
Machine A                    Machine B
┌─────────────┐             ┌─────────────┐
│ Local Index  │◄──gRPC──►│ Local Index  │
│ +            │             │ +            │
│ Remote Proxy │             │ Remote Proxy │
└─────────────┘             └─────────────┘
       │                           │
       ▼                           ▼
┌─────────────┐             ┌─────────────┐
│ Qdrant      │             │ Qdrant      │
│ (local shard)│             │ (local shard)│
└─────────────┘             └─────────────┘
```

---

## Phase 8 — Ecosystem & API (Weeks 22-24)
**Theme:** ASH is a platform. Ship the API, the plugin system, and the community tools.

### Decision Fork: Plugin Runtime
**Chosen:** Shell scripts + Python plugins with a JSON API contract. No WASM runtime (too early). Plugins register via `.plugin` files in `/etc/ash/plugins/`.

### Deliverables
- [x] Query handler plugin system (FUTURE #69)
- [x] Webhook notification support (FUTURE #34)
- [x] `ash-ask` — RAG question answering API
- [x] GitHub Actions integration — `ash-ci` runs in CI
- [x] VSCode extension — search from editor
- [x] Neovim plugin — `:AshSearch` command
- [x] API documentation — OpenAPI spec for all services
- [x] Example plugins — GitHub, Notion, Obsidian, Jira
- [x] Community plugin registry (`ash plugin search`)

### API Surface (v1)
```
POST   /api/search          — semantic + keyword hybrid
POST   /api/ask             — RAG question answering
GET    /api/status          — service health
GET    /api/stats           — storage and performance metrics
POST   /api/index/scan      — trigger reindex
POST   /api/config          — update config (live reload)
DELETE /api/workspace/:id   — destroy disposable workspace
```

---

## Dependency Map

```
Phase 1 ───► Phase 2 ───► Phase 3 ───► Phase 8
   │             │            │
   ▼             ▼            ▼
Phase 4 ───► Phase 5 ───► Phase 6 ───► Phase 7
```

Phases 1-3 are sequential (core stack). Phases 4-5 can partially overlap with 3. Phase 6-7 should be sequential (need observability before federation). Phase 8 can start after Phase 3.

---

## Priority Sequencing

| Priority | Phase | Feature | Effort | Impact |
|:--------:|:-----:|---------|:------:|:------:|
| P0 | 2 | BM25 keyword fallback | 2 days | Unlocks exact-match search |
| P0 | 2 | Inotify file watcher | 1 day | Fixes 60s poll latency |
| P0 | 3 | `ash workspace` commands | 3 days | Core vibecoding UX |
| P0 | 4 | Cross-distro support | 5 days | Unlocks 80% of Linux market |
| P1 | 1 | Offline bundle | 2 days | Air-gapped deploys |
| P1 | 3 | `ash-ask` RAG | 3 days | Killer feature |
| P1 | 5 | AppArmor profiles | 2 days | Security baseline |
| P1 | 6 | `ash-tui` dashboard | 5 days | Visibility |
| P2 | 7 | Federation | 10 days | Advanced |
| P2 | 8 | Plugin system | 8 days | Ecosystem |

---

## How This Document Evolves

1. **Every month**: phases are re-prioritized based on user feedback and GitHub issues
2. **Every commit**: update the status of deliverables (✅ 🔧 📅)
3. **Every release**: move completed phases to `CHANGELOG.md`, update the dependency map
4. **The `agy` tool reads this file** — `agy status` shows phase completion, `agy next` shows what to work on

---

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|:----------:|:------:|------------|
| Cross-distro package drift | High | Medium | CI tests all distros; community PRs for fixes |
| Go rewrite stalls | Medium | High | Strangler fig — ship Go CLIs one at a time; bash stays as fallback |
| Container mode perf gap | Medium | High | Profile before optimizing; document trade-offs |
| Plugin API changes after v1 | Low | High | Mark v1 as experimental; stabilize by v5.0 |
| LLM model deprecation | Medium | Medium | Support `--model` flag; nomic-embed-text is stable but have fallback list |

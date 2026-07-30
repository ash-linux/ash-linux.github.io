# FUTURE — 100x Improvements for Ash Installer

> 70 high-impact features and improvements for `scripts/ash-install.sh`, organized by theme.

---

## 🛡️ Reliability & Safety

### 1. Dry-run / preview mode
Add `--dry-run` flag that shows every command without executing. Users can inspect exactly what'll change before the install touches their system. Each phase prints: "Would install: curl, wofi, qdrant v1.13.6 ...".

### 2. Network resilience with exponential retry backoff
Every `curl|bash` and GitHub API call uses exponential backoff (1s, 2s, 4s, 8s, 16s) and a configurable `--max-retries` flag. Currently, a single network blip aborts the entire install. Use:
```bash
retry() { local n=0; until "$@" || ((++n == max)); do sleep $((2**n)); done; }
```

### 3. Checksum verification of every downloaded binary
Download SHA256 checksums from GitHub releases alongside Qdrant/Ollama binaries and verify before installing. Abort on mismatch with full error. Currently, any compromised mirror injects malicious binaries without detection.

### 4. Atomic install with commit/rollback transactions
Wrap the entire install in a `systemd-run --scope --unit=ash-install` so if the machine crashes mid-install, systemd reports the unit as failed on next boot and triggers automatic rollback from the btrfs snapshot.

### 5. Full state-machine resume (not just skip-if-running)
Save granular state per-substep (not per-phase). On re-run, resume from the exact failed substep rather than re-doing the whole phase. E.g., if Qdrant binary downloads but service fails, resume from `systemctl enable`. Use a JSON state file:
```json
{"phase":"qdrant","step":"download","status":"done","checksum":"abc123"}
```

### 6. Interactive rollback menu
After a crash, offer a TUI menu with: (a) restore snapshot + retry, (b) restore snapshot + abort, (c) continue despite failure, (d) show full diff of changes made. Currently rollback is silent and automatic which can be destructive if the user had unrelated changes.

### 7. Preflight system compatibility check
Before any changes, run a comprehensive check and display a report:
- Kernel ≥ 6.x (for btrfs, io_uring)
- btrfs-progs available
- Disk space ≥ 10GB free
- RAM ≥ 4GB (8GB recommended)
- Python ≥ 3.10
- systemd ≥ 250
- `sudo` configured without tty requirement
- no conflicting services on ports 6333, 11434
Fail early with clear messages instead of mid-install cryptic errors.

### 8. Power-failure resilience
Add `sync` before each critical phase. Use `flock` on the state file to prevent concurrent installs. Detect incomplete writes via checksum markers at the end of the state file.

### 9. Backup of every modified config file
Before touching `hyprland.conf`, `.bashrc`, `.zshrc`, `/etc/systemd/...`, copy the original to `/var/backups/ash-install/` with a timestamp. `ash-doctor --restore-configs` can revert any single file.

### 10. SELinux / AppArmor profile generation
Ship AppArmor profiles for qdrant, ollama, and lsfs-daemon that confine each to only the resources they need. Enforce during `verify_all`. Currently, all three run unrestricted.

---

## 🚀 Performance & Resource Management

### 11. GPU detection + model offloading
Auto-detect NVIDIA (nvidia-smi), AMD (rocm-smi), or Intel (intel_gpu_top) GPU, then configure Ollama to offload embedding model layers to GPU. Add `--gpu nvidia|amd|intel|cpu` override. For NVIDIA, set `OLLAMA_NUM_GPU=1` and `OLLAMA_GPU_OVERHEAD=512`. Currently everything runs on CPU, making the embedding model 10x slower than necessary on GPU-capable hardware.

### 12. Configurable resource limits
Add `--memory-max=2G`, `--cpu-quota=50%`, `--disk-cache=1G` flags. Generate systemd unit drop-ins (`/etc/systemd/system/qdrant.service.d/limits.conf`) that apply these limits. Critical for running on low-spec VMs or alongside other workloads.

### 13. Incremental LSFS indexing with inotify
Replace the 60-second polling loop with `inotify` (via `watchdog` Python library or `inotifywait` for the initial scan). Index files within milliseconds of creation/modification. Drop CPU usage from "constantly polling" to zero.

### 14. Binary caching layer
Cache downloaded Qdrant/Ollama binaries to `/var/cache/ash-install/`. Use content-addressable storage (SHA256 filenames). On re-run or repair, skip re-download if checksums match. Set up a `systemd-tmpfiles` rule to auto-clean after 30 days.

### 15. Parallel model pull with real-time progress
Pull `nomic-embed-text` using Ollama's streaming API and show a live progress bar (bytes downloaded / total, speed, ETA). For multiple models, pull in parallel with separate progress bars in the TUI.

### 16. Optimized file walker with compiled .lsfsignore
Parse `.lsfsignore` once into a compiled regex trie instead of iterating through IGNORE list for every single file. Use `pathspec` library or implement a `fnmatch` tree walker. On a home dir with 500k files, this is the difference between 2 minutes and 20 minutes.

### 17. Memory-mapped Qdrant storage configuration
Configure Qdrant with `--optimizer-overcommit` and `--mmap-threshold-kb=10000` on systems with >8GB RAM. Detects total RAM and adjusts Qdrant's `optimizer_cpu_limit`, `segment_number_threshold`, and HNSW `ef_construct` based on available cores.

### 18. JIT compilation of indexing pipeline
Use `functools.lru_cache` on the embedding call and batch embedding requests (send 16 text chunks at once via Ollama's `/api/embed` batch endpoint instead of 16 sequential `/api/embeddings` calls). Reduces model inference overhead by 10-15x during initial full scan.

### 19. Swap-aware install
Detect if swap is enabled. If not and RAM < 8GB, offer to create a 4GB swapfile. Without swap, Ollama + Qdrant together can OOM on 4GB machines during concurrent model loading.

### 20. Build from source fallback
If prebuilt binaries fail checksum or are unavailable for the target architecture (ARM64, RISC-V), fall back to building Qdrant from source (via `cargo`) and Ollama from source (via `go`). Show estimated build time in the TUI.

---

## 🧠 AI & Search Capabilities

### 21. Multi-model search (semantic + BM25 hybrid)
Add a BM25 keyword-search fallback via SQLite FTS5 or `tantivy` that merges with vector search results using Reciprocal Rank Fusion. "Exact-match" queries like "auth.go" or "Dockerfile" work even when the embedding model misses them. The fusion weight is tunable: `--search-mode hybrid|semantic|keyword`.

### 22. Query categorization & routing
Classify queries into types at the launcher level using a small on-device classifier (or simple regex/pattern matching):
- `[file]` → LSFS semantic search
- `[app]` → desktop-file search + gtk-launch
- `[calc]` → Python eval sandbox
- `[web]` → open in browser
- `[cmd]` → shell execution with confirmation
- `[clip]` → clipboard history search
One search bar to rule them all. Show a type badge next to each result.

### 23. Reindex scheduler with debouncing
Batch rapid file events (e.g., `git clone` creates 10k files) using a debounce timer: only trigger indexing after 30s of filesystem quiet. Store a fast hash map of `(inode, mtime) -> etag` to skip unmodified files entirely rather than re-embedding every time.

### 24. Multi-user install
Detect all human users (UID ≥ 1000, shell != nologin) from `/etc/passwd`. Prompt to install for each selected user. Run per-user LSFS daemons under separate systemd user services. Store shared indexes in `/var/lib/ash/shared` and per-user indexes in `~/.ash/index`. Each user's Super+Space is isolated.

### 25. Embedding model selection at install time
TUI menu to choose embedding model with live size/speed/quality comparison:

| Model | Size | Quality | Speed |
|-------|------|---------|-------|
| nomic-embed-text | 137MB | Good | Fast |
| mxbai-embed-large | 334MB | Better | Medium |
| snowflake-arctic-embed2 | 1.2GB | Best | Slow |
| all-MiniLM-L6-v2 | 80MB | Fair | Fastest |

### 26. Content extraction for binary files
Add `python-magic`, `textract`, `pdftotext`, `tesseract`, `pandoc`, and `ffmpeg` integration so these file types get indexed:
- PDF → text via pdftotext
- DOCX/PPTX/XLSX → text via python-docx/python-pptx/openpyxl
- Images → OCR via tesseract
- Audio → transcription via whisper.cpp (optional, large dep)
- Archives → listing of contents

### 27. Hybrid local + remote model support
Allow Ollama to pull a small local model for fast queries (nomic-embed-text) and route to a remote OpenAI/Anthropic API for complex understanding queries. Store the API key in the system keyring (`secret-tool` / `keyctl`). Add `--remote-provider openai|anthropic|generic --remote-model gpt-4o`.

### 28. File relationship graph
Build a knowledge graph of file relationships during indexing: imports, symlinks, git history, recent co-access patterns. Store in Qdrant as hybrid vector+graph index. Query "what depends on this config" or "files I was working on with that script".

### 29. Semantic clipboard history
Index clipboard content in real-time via `wl-clipboard` watch mode. Press `Super+Space` then type "the email I copied yesterday" to find it. Store clipboard entries with timestamps and embedding vectors in a circular buffer (last 1000 entries).

### 30. RAG over personal files
Add a `ash-ask` command that given a natural-language question:
1. Vector-search all indexed files for relevant chunks
2. Feed chunks + question to Ollama with a QA prompt
3. Return a synthesized answer with source citations
"Ash, what's the timezone configured in my app?" → reads config files and answers.

---

## 📊 Observability & UX

### 31. Real-time install dashboard with split-pane TUI
Show a live-updating dashboard using `tput` cursor movement:
```
┌─ Phases ────────────────────┬─ Current Output ─────┐
│ ✔ System Packages           │ Installing curl...   │
│ ✔ Qdrant Vector DB          │ Downloading 45% ████ │
│ ◌ Ollama AI Engine          │                      │
│ ◌ Embedding Model           │                      │
│ ◌ LSFS Indexer              │                      │
│ ◌ Desktop Launcher          │                      │
│ ◌ Auto Updates              │                      │
│ ─────────────────────────── │                      │
│ Overall: 25% ███████░░░░░░░ │                      │
├─────────────────────────────┴──────────────────────┤
│ ⚠ Ollama: first pull can take 3-10 min             │
└────────────────────────────────────────────────────┘
```

### 32. `ash-doctor` comprehensive diagnostic tool
Ship a standalone diagnostic script (`/usr/local/bin/ash-doctor`) that runs 40+ checks and exports to a shareable report:
- Service health (systemd, ports, process existence)
- Log tail (last 50 lines of each service)
- Python module imports
- Disk space alerts
- GPU availability
- File permissions
- SELinux/AppArmor status
- DNS resolution
- NTP sync status
- Btrfs health (`btrfs device stats`)
- Snapshot list
- LSFS index count and age
Output: color terminal + `/var/log/ash-doctor-report.txt` + optional JSON.

### 33. `ash-uninstall` comprehensive uninstaller
Generate a manifest of every file created during install (`/var/log/ash-install/manifest.txt`). Ship `ash-uninstall` that:
1. Stops all services (systemd user + system)
2. Removes all installed binaries
3. Removes systemd units and timers
4. Purges Qdrant vector storage
5. Removes Ollama models and binary
6. Restores backed-up config files
7. Removes user services and PATH additions
8. Uninstalls packages that were installed and *not* pre-existing

### 34. `--webhook` notification support
Optional Slack/Discord/ntfy/email webhook URL. Gets notified on:
- Install completed successfully
- Install failed (with phase name + log snippet)
- Auto-update ran and result
- ash-doctor detects a degraded state
Configurable via `--webhook-url https://hooks.slack.com/...` or `/etc/ash/webhook.conf`.

### 35. HTML install report with waterfall chart
Instead of just terminal output, write `/var/log/ash-install/report.html` with:
- Waterfall timeline chart (install phases, their duration, parallelism)
- Log viewer (searchable, expandable sections per phase)
- System info snapshot (CPU, RAM, disk, kernel, GPU)
- Error callouts with suggested fixes
- Machine-readable JSON export alongside

### 36. Automatic log rotation and archiving
Ship a `logrotate.d` config for `/var/log/ash-install/*.log`:
- Rotate weekly
- Keep 4 weeks
- Compress with zstd
- Optionally upload to S3-compatible storage on rotation

### 37. `Super+Space` sound effect
Add an optional audio cue (a short `.opus` or `.wav`) on launcher open/close. Play via `pw-play` or `paplay`. Configurable: `--sound-theme minimal|retro|none`.

### 38. Welcome wizard on first run
After install, instead of just a dashboard, show a 3-step TUI wizard:
1. "Press Super+Space to test" — walks through a sample search
2. "Customize your models" — choose default embedding model
3. "Enable telemetry?" — opt-in with privacy notice

### 39. First-run onboarding results
On first Super+Space press, run a welcome query: "Your files are being indexed. Here's what that means..." with pre-indexed Ash documentation so the user sees instant results instead of an empty screen.

### 40. Install telemetry (opt-in)
After install, optionally POST anonymous stats (SHA256 of machine ID, phases + durations, distro, kernel, errors, hardware class). All displayed upfront with full transparency. `--no-telemetry` to disable. Public dashboard at `ash.sh/status` showing aggregate install stats.

### 41. Desktop notifications during install
Use `notify-send` during long-running phases: "Qdrant: downloading (45%)", "Embedding model: pulling (5/10 min remaining)". Keeps the user informed even if they alt-tab away from the terminal.

### 42. Dark/light theme detection for TUI
Detect terminal theme via `osascript -e '...'` (macOS) or `gsettings` (GNOME) or `kitty @ ls` and adjust TUI colors accordingly. `--color-scheme auto|dark|light`.

---

## 🔧 Cross-Platform & DevOps

### 43. Non-Arch Linux support (Debian/Ubuntu/Fedora/openSUSE)
Detect package manager (`apt`, `dnf`, `zypper`, `apk`, `xbps`) and translate package names via a lookup table. Qdrant/Ollama binaries are distro-agnostic. This single change unlocks 80%+ of the Linux market.

### 44. Docker/podman rootless install mode
Add `--container` flag that installs everything (Qdrant, Ollama, LSFS daemon) as rootless podman/docker quadlets or `podman-compose` stacks instead of bare-metal systemd services. Enables Ash on:
- Fedora Silverblue / Universal Blue
- SteamOS
- NixOS
- Container-optimized OSes

### 45. Offline/air-gapped install bundle
`--offline /path/to/bundle.tar.gz` mode where the bundle contains all binaries, debs/rpms, and models with checksums. Generate the bundle with `ash-install --bundle`. Resumable download of the bundle with aria2c.

### 46. Config profiles: desktop / server / ci / minimal
- `--profile desktop` (default) — full stack with Hyprland, wofi, Super+Space binding
- `--profile server` — Qdrant + Ollama + lsfs-query CLI only, no GUI deps
- `--profile ci` — verification tooling only (ash-doctor + health checks)
- `--profile minimal` — LSFS CLI only, no Qdrant (uses SQLite FTS5 instead)
Each profile = different set of phases. Skips irrelevant downloads and configs.

### 47. Self-update with changelog diff
Before install, check GitHub releases API for newer version. If running `v2.5.0` and `v3.0.0` exists, show a semantic diff changelog and prompt to upgrade. Post-install, store the installed script hash in the state file so `ash-doctor --verify-script` can detect tampering.

### 48. macOS support via Homebrew
Detect macOS, use Homebrew for dependencies, launch Qdrant/Ollama as launchd services instead of systemd. LSFS indexes ~/Documents, ~/Desktop, ~/Downloads. Super+Space opens via a Raycast/Alfred extension or global hotkey via `skhd`.

### 49. Windows WSL2 support
Detect WSL2 environment, skip Linux-specific phases (Hyprland, systemd), use Windows interop for notifications, and integrate with Windows Search index. Register as a WSLg application so it appears in the Start Menu.

### 50. Ansible role generation
After an interactive install, optionally output an Ansible role or Terraform module that reproduces the same configuration non-interactively. `--export-ansible /path/to/role/`. This makes Ash deployable at enterprise scale.

### 51. Vagrant box provisioning
Integrate with the existing `packer/vagrant-box.pkr.hcl` — after building an Ash ISO, auto-install via `ash-install --vagrant` that detects it's inside a Vagrant build and skips the TUI/systemd steps that don't apply.

### 52. Ephemeral install for CI/testing
`--ephemeral` flag that runs the full install, runs `ash-doctor`, then runs `ash-uninstall`. Return exit code 0 only if all phases pass. Designed for CI pipelines (GitHub Actions, GitLab CI) to validate every commit.

### 53. Multi-architecture binary selection
Detect architecture (`uname -m`) and download appropriate Qdrant/Ollama build:
- x86_64 → x86_64-unknown-linux-musl
- aarch64 → aarch64-unknown-linux-musl
- armv7l → armv7-unknown-linux-musleabihf
- riscv64 → riscv64gc-unknown-linux-gnu
Currently hardcoded to x86_64.

### 54. Flatpak / Snap / AppImage integration
Index Flatpak, Snap, and AppImage applications so they appear in `[app]` search results. Parse `.desktop` files from `/var/lib/flatpak/exports/share/applications/` and `/snap/bin/`. Open via `flatpak run` or `snap run` respectively.

---

## 🔐 Security & Hardening

### 55. Systemd service hardening
Apply systemd security features to all shipped units:
```
ProtectSystem=strict
ProtectHome=true
PrivateTmp=true
NoNewPrivileges=true
CapabilityBoundingSet=~CAP_SYS_ADMIN
MemoryDenyWriteExecute=true
RestrictNamespaces=true
```
Override-able via drop-ins. Currently all three services run with default (unrestricted) permissions.

### 56. API authentication for Qdrant
Generate a random API key on install, configure Qdrant with `QDRANT__SERVICE__API_KEY`, and store in `/etc/ash/secrets.env` (0600 permissions). LSFS and query tools read from the secrets file. Prevents local privilege escalation via localhost APIs.

### 57. Firewall rules
Generate iptables/nftables rules to restrict Qdrant and Ollama to localhost only:
```
nft add rule inet filter input ip daddr 127.0.0.1 tcp dport { 6333, 11434 } accept
nft add rule inet filter input tcp dport { 6333, 11434 } drop
```

### 58. Audit log
Log every search query to a append-only audit log at `/var/log/ash/audit.log` with timestamp, truncated query, result count, and latency. Use `chattr +a` to make the log append-only. Configurable: `--audit-level none|metadata|full`.

### 59. Verified boot measurement
Take a TPM PCR measurement after installation completes. On subsequent boots, `ash-doctor --verify-boot` checks that the measurement matches, detecting tampering of installed binaries. Requires `tpm2-tools`.

### 60. GPG signing of install manifest
Sign the install manifest (`/etc/ash-installed`) with a generated GPG key. Future `ash-doctor --verify` checks the signature. If the manifest is tampered, the user is alerted.

### 61. Secrets encryption
Store API keys, tokens, and any collected telemetry identifiers encrypted at rest using `age` encryption with a key derived from the machine's TPM/systemd credentials.

### 62. USB/offline signing
For air-gapped deploys, support detached GPG signatures: `ash-install --sign` GPG-signs the bundle, `ash-install --verify-bundle bundle.tar.gz.asc` verifies before installing.

---

## 🗄️ Storage & Data Management

### 63. Qdrant storage tiering
Configure Qdrant with hot/warm/cold tiers:
- Hot: in-memory segments (fastest, 1GB max)
- Warm: mmap'd segments on NVMe (default)
- Cold: compressed segments on HDD (for files >30 days old)
Auto-migrates points between tiers based on access recency using Qdrant's built-in optimizer.

### 64. Index deduplication and compression
Before storing embeddings, deduplicate identical file content (by SHA256). Store one embedding per unique content hash with a reference count. Reduces storage by 30-50% for projects with many copies/config variants.

### 65. Database backup schedule
Ship a `ash-backup.timer` that:
1. Snapshot Qdrant's storage directory via btrfs subvolume
2. Dump Ollama model list
3. Archive to `/var/backups/ash/`
4. (Optional) Push to S3/B2/rsync.net
Configurable retention: `--backup-retain 7`.

### 66. Btrfs auto-defrag for Qdrant storage
Set up a `btrfs-auto-defrag` path for `/var/lib/qdrant/storage` and `/var/lib/ollama`. Qdrant's write-heavy workload fragments btrfs over time; periodic defrag keeps performance consistent.

### 67. Storage usage dashboard
Add `ash-stats` command that shows:
```
Storage:
  Index size: 2.3 GB (1,234 files)
  Qdrant segments: 12 (4 memory, 8 disk)
  Ollama models: 2 (1.4 GB)
  Cache: 340 MB
  Snapshot count: 3 (using 2.1 GB)
```

### 68. Remote index federation
Optionally connect multiple Ash machines so they share a unified index. A query on machine A returns results from machines B and C (with network latency compensation). Uses Qdrant's built-in raft consensus or a lightweight gRPC proxy.

---

## 🤖 Developer & Power User

### 69. Plugin system for query handlers
Design a plugin API so third-party handlers can register with the launcher:
```bash
# /etc/ash/plugins/github.conf
[handler]
match = ^pr|^issue|^repo
command = /usr/local/bin/ash-plugin-github "$QUERY"
icon = github
```
Plugins receive the query, return JSON results, and can open files/apps/URLs. Ship with example plugins for GitHub, Notion, Obsidian, and Jira.

### 70. Live reload of configuration
All Ash configuration lives under `/etc/ash/` and supports `SIGHUP` reload:
- `pkill -HUP lsfs-daemon` → reloads `.lsfsignore` and config
- Qdrant config changes → `systemctl reload qdrant`
- No restart needed for config changes to take effect
<![CDATA[<div align="center">

# 🌿 Ash Linux

### The AI-Native Operating System

**Search your files by meaning. Run AI locally. No cloud. No API keys. Just press `Super+Space`.**

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Arch Linux](https://img.shields.io/badge/Based%20on-Arch%20Linux-1793D1?logo=archlinux&logoColor=white)](https://archlinux.org)
[![Platforms](https://img.shields.io/badge/Platforms-x86__64%20%7C%20ARM64-blue)]()
[![Offline](https://img.shields.io/badge/Works-100%25%20Offline-brightgreen)]()

---

**An Arch Linux environment with built-in semantic file search, vector memory, and an AI-native desktop launcher.**
**Boot it up, press a key, and find any file by *what it contains* — not just its name.**

[Get Started](#-get-started) · [Features](#-features) · [How It Works](#-how-it-works) · [FAQ](#-frequently-asked-questions) · [Docs](docs/)

</div>

---

## ✨ What is Ash?

Ash is an **AI-powered Linux desktop** that understands your files. Instead of remembering filenames and folder paths, you describe what you're looking for in plain English — and Ash finds it instantly.

> **Think of it like Spotlight on macOS, but it actually understands your files' *content and meaning*.**

- 🔍 **"Find my config files"** → returns all configuration files, regardless of name or location
- 📝 **"That script I wrote for backups"** → finds it by understanding what the code does
- ⏰ **"Files from the last 2 hours"** → time-based search as a fallback
- 🔒 **Everything runs locally** — your files never leave your machine

---

## 🚀 Get Started

### Option 1: One-Line Install (on existing Arch Linux)

Already running Arch? Deploy Ash on your existing system in under 2 minutes:

```bash
curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash
```

That's it. Reboot, press **Super+Space**, and start searching.

### Option 2: Download the ISO

Get the full Ash experience — a pre-built Arch Linux image with everything pre-configured:

| Format | Platform | Download |
|--------|----------|----------|
| **ISO** (x86_64) | VMware / VirtualBox / bare metal | [Latest Release](https://github.com/exonew2/files/releases) |
| **ISO** (ARM64) | Apple Silicon / RPi / ARM servers | [Latest Release](https://github.com/exonew2/files/releases) |
| **OVA** | VMware Workstation / Fusion | [Latest Release](https://github.com/exonew2/files/releases) |
| **QCOW2** | QEMU / KVM / Proxmox | [Latest Release](https://github.com/exonew2/files/releases) |
| **Vagrant Box** | Vagrant + VirtualBox / libvirt | [Latest Release](https://github.com/exonew2/files/releases) |
| **WSL2 Tarball** | Windows Subsystem for Linux | [Latest Release](https://github.com/exonew2/files/releases) |
| **RPi Image** | Raspberry Pi 4 / 5 | [Latest Release](https://github.com/exonew2/files/releases) |

### Option 3: Docker / Container

```bash
docker pull ghcr.io/ash-linux/ash:latest
```

Also available on [Docker Hub](https://hub.docker.com/r/ashlinux/ash) and [Quay.io](https://quay.io/repository/ash-linux/ash).

---

## 🧠 Features

### Semantic File Search — `Super+Space`

Press **Super+Space** anywhere on the desktop. A search bar appears. Type a natural-language query:

| You type | What Ash finds |
|----------|----------------|
| `config` | All configuration files — `.conf`, `.yaml`, `.toml`, scripts with config logic |
| `backup script` | Shell scripts related to backups, regardless of filename |
| `TODO` | Files containing TODO items, task lists, notes |
| `files from 3 days` | Everything modified in the last 3 days (time-based fallback) |
| `docker setup` | Dockerfiles, compose files, container-related scripts |

Results appear in a clean selection menu. Pick a file and it opens instantly:
- **Code/text files** → open in Neovim (in Kitty terminal)
- **Directories** → open in Yazi file manager
- **Applications** → launch directly

### 🤖 Fully Local AI Stack

Everything runs on your machine. No internet required after initial setup.

| Component | What It Does |
|-----------|-------------|
| **Ollama** | Runs the AI embedding model locally |
| **nomic-embed-text** | Converts your files into meaning-vectors (768 dimensions) |
| **Qdrant** | Vector database — stores and searches file embeddings |
| **LSFS Daemon** | Watches your filesystem and auto-indexes new/changed files |

### 🖥️ Beautiful Desktop Environment

- **Hyprland** — Modern Wayland tiling compositor
- **Catppuccin Mocha** — Elegant dark theme (purple/navy palette, pink accents)
- **Kitty** — GPU-accelerated terminal
- **Wofi** — Clean application launcher
- **Auto-login** — Boots straight to the desktop, no login screen

### ⚡ Zero Configuration

- **Auto-login** → Boots directly into Hyprland desktop
- **Auto-start** → All AI services launch at boot
- **Auto-index** → New files are automatically embedded and searchable
- **Idempotent install** → Re-run the installer anytime — it never breaks things

### 🔒 Privacy First

- **100% offline** — No data leaves your machine, ever
- **No API keys** — No accounts, no subscriptions, no telemetry
- **All services on localhost** — Ollama on `127.0.0.1:11434`, Qdrant on `127.0.0.1:6333`
- **VM isolation** — Run in a virtual machine for additional sandboxing
- **Systemd sandboxing** — Services run with CPU/memory limits and restricted permissions

---

## 🔧 How It Works

```
You press Super+Space
        │
        ▼
┌─────────────────────┐
│   Search Prompt     │  ← Type any query: "backup scripts", "TODO", etc.
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│   Ollama API        │  ← Converts your query into a 768-dim vector
│   nomic-embed-text  │     (runs locally, no internet needed)
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│   Qdrant Database   │  ← Finds the most similar files by meaning
│   Vector Search     │     (cosine similarity on all indexed files)
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│   Results Menu      │  ← Select a file to open it
│   wofi              │
└─────────────────────┘
```

**Behind the scenes**, the LSFS daemon continuously watches your filesystem. When you create or edit a file, it automatically:
1. Reads the file content
2. Sends it to Ollama for embedding
3. Stores the vector in Qdrant

This means every file is searchable by meaning within seconds of being saved.

---

## 📦 What's Included

### Desktop & Shell
| Tool | Description |
|------|-------------|
| Hyprland | Wayland tiling compositor with Catppuccin Mocha theme |
| Kitty | GPU-accelerated terminal emulator |
| Wofi | Application launcher and search UI |
| Neovim | Text editor (opens from search results) |
| Yazi | Terminal file manager |
| Starship | Modern shell prompt |
| Zsh + Fish | Modern shells, pre-configured |
| ripgrep, fd, bat, eza, fzf | Modern CLI replacements for grep, find, cat, ls |

### AI & Search
| Tool | Description |
|------|-------------|
| Ollama | Local LLM and embedding model server |
| nomic-embed-text | 768-dim embedding model for semantic search |
| Qdrant | Vector database for file embeddings |
| LSFS Daemon | Auto-indexes filesystem changes in real-time |
| Launcher Hook | Pure-bash semantic search in ~50 lines |

### Development
| Tool | Description |
|------|-------------|
| Python, Node.js, Go, Rust | Full development stacks |
| Docker + Podman | Container runtimes |
| Terraform + Ansible | Infrastructure as code |
| Git + GitHub CLI | Version control |

### AI Agent Integration
| Tool | Description |
|------|-------------|
| Claude Desktop | Pre-configured MCP server connection |
| Cursor | Rules file with project context |
| Windsurf | Rules file with project context |
| Cline | MCP config for VS Code |
| Codebase Memory | Knowledge graph for persistent AI agent context |

---

## 🖥️ Platform Support

### Recommended: VMware

Ash is optimized for VMware Workstation / Fusion. After importing:

1. **Enable 3D acceleration** in VM Settings → Display
2. **Add to your `.vmx` file** (on the host machine):
   ```
   mks.enableVulkanRenderer = "FALSE"
   svga.disableFIFO = "TRUE"
   ```
3. **For clipboard** (optional):
   ```
   isolation.tools.copy.disable = "FALSE"
   isolation.tools.paste.disable = "FALSE"
   isolation.tools.setGUIOptions.enable = "TRUE"
   ```

See [VMware Setup Guide](docs/vmware-setup.md) for details.

### Also Works On

| Platform | Notes |
|----------|-------|
| **VirtualBox** | Guest additions included |
| **QEMU/KVM** | SPICE agent included, QCOW2 format available |
| **Proxmox** | Import QCOW2 image |
| **Bare Metal** | Install from ISO, full GPU support |
| **Raspberry Pi 4/5** | ARM64 image available |
| **WSL2** | Tarball available, no GUI (CLI search only) |
| **AWS / GCP / Azure** | Cloud images available via Packer + Terraform |

---

## 🔄 Updating

Re-run the installer. It's idempotent — it updates everything without destroying your data:

```bash
curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash
```

Your Qdrant vector database and Ollama models are preserved across updates.

For individual component updates, see the [Update Guide](docs/updating.md).

---

## 🏗️ Multi-Node Cluster (Advanced)

Scale Ash across multiple machines with the built-in cluster controller:

```bash
# Spin up a 3-node cluster with load-balanced Ollama
ash-cluster up --nodes 3

# Check status
ash-cluster status

# Tear down
ash-cluster down
```

The cluster uses Docker Compose or HashiCorp Nomad, with Consul for service discovery and Traefik for load balancing. See [Cluster Guide](docs/cluster.md) for details.

---

## ❓ Frequently Asked Questions

<details>
<summary><b>How much disk space does Ash need?</b></summary>

The base system needs ~5 GB. The AI model (`nomic-embed-text`) adds ~300 MB. The Qdrant vector database grows based on how many files you have indexed — typical usage is under 500 MB.
</details>

<details>
<summary><b>How much RAM does it need?</b></summary>

Minimum 4 GB, recommended 8 GB. Ollama uses ~1 GB for the embedding model. Qdrant is limited to 2 GB by systemd configuration.
</details>

<details>
<summary><b>Does it need a GPU?</b></summary>

No. The `nomic-embed-text` model runs comfortably on CPU. A GPU provides modest speedup for batch embedding but is not required.
</details>

<details>
<summary><b>Does it need internet?</b></summary>

Only for the initial setup (downloading packages, the AI model, and Qdrant binary). After that, everything runs 100% offline.
</details>

<details>
<summary><b>What files does it index?</b></summary>

By default, it indexes `~/.config/scripts` and configurable paths. You can control what gets indexed by editing the `~/.lsfsignore` file (uses gitignore-compatible syntax). Binary files, images, and common junk directories are excluded by default.
</details>

<details>
<summary><b>Is my data safe?</b></summary>

Yes. Everything runs on localhost only. No ports are exposed to the network. No data is sent anywhere. The VM boundary provides additional isolation from your host machine.
</details>

<details>
<summary><b>Can I use this on a non-Arch distro?</b></summary>

The core concept (Ollama + Qdrant + bash launcher) can work on any Linux distro, but the automated installer targets Arch Linux specifically. Porting the script to Ubuntu/Fedora would require replacing `pacman` calls with `apt`/`dnf`.
</details>

<details>
<summary><b>How is this different from GNOME Search / Spotlight?</b></summary>

Traditional search matches filenames and file contents literally. Ash uses **vector embeddings** — it understands the *meaning* of your query and finds semantically similar files. Searching "backup" will find a script called `snapshot.sh` if it contains backup logic, even though the word "backup" never appears in the filename.
</details>

---

## 📖 Documentation

| Guide | Description |
|-------|-------------|
| [Quick Start](docs/quickstart.md) | Install on existing Arch Linux in 2 minutes |
| [VMware Setup](docs/vmware-setup.md) | Optimizing Ash for VMware Workstation/Fusion |
| [How Search Works](docs/how-search-works.md) | Deep dive into the semantic search pipeline |
| [Updating](docs/updating.md) | Keeping Ash and its components up to date |
| [Backup & Restore](docs/backup-and-restore.md) | Protecting your data and disaster recovery |
| [GPU & Performance](docs/gpu-and-performance.md) | GPU passthrough and performance tuning |
| [Cluster Mode](docs/cluster.md) | Multi-node AI cluster with load balancing |
| [AI Agent Integration](docs/ai-agents.md) | Connecting Claude, Cursor, Windsurf, and Cline |
| [Security](docs/security.md) | Security model and threat surface |
| [Comparison](docs/comparison.md) | How Ash compares to other AI OS projects |
| [Health Checks](docs/health-checks.md) | Verifying everything is running correctly |
| [Troubleshooting](docs/troubleshooting.md) | Common issues and fixes |

---

## 🗺️ Roadmap

- [ ] GUI settings panel for configuring indexed paths
- [ ] Multi-model support (switch embedding models on the fly)
- [ ] Qdrant backup/restore automation
- [ ] Image and PDF content indexing
- [ ] Firewall configuration out of the box
- [ ] AppArmor security profiles
- [ ] Ubuntu/Fedora installer variants

---

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**Areas where help is most needed:**
- VMware clipboard automation (host-side scripting)
- Additional embedding model support
- Qdrant backup/restore tooling
- Ports to Ubuntu and Fedora

---

## 📄 License

MIT License — see [LICENSE](LICENSE) for details.

---

<div align="center">

**Built with ❤️ on Arch Linux**

*Ash Linux — Your files, understood.*

</div>
]]>

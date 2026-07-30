<![CDATA[# Quick Start

Get Ash Linux running on your existing Arch Linux system in under 2 minutes.

## What You Need

- **Arch Linux** — a running installation (VM or bare metal)
- **Internet connection** — for the initial download only
- **sudo access** — the installer needs root to install packages and services

That's it. No pre-installed dependencies required.

## Install

Run this single command:

```bash
curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash
```

The installer is **idempotent** — you can safely run it again at any time. It will never break your existing setup.

### What the Installer Does

1. **Tunes your system** — Sets inotify limits, swappiness, and other sysctl parameters for optimal performance
2. **Installs packages** — `wofi`, `jq`, and other dependencies via `pacman`
3. **Downloads Qdrant** — Standalone binary from GitHub releases (no Docker, no AUR)
4. **Configures Ollama** — Enables the systemd service, pulls the `nomic-embed-text` embedding model
5. **Deploys LSFS** — Installs the launcher hook and indexing daemon to `~/.config/scripts/`
6. **Sets up services** — Creates systemd units, enables auto-start on boot
7. **Configures Hyprland** — Binds `Super+Space` to the semantic search launcher
8. **Enables auto-login** — Boots straight into Hyprland desktop

## First Use

1. **Reboot** your system (or restart Hyprland with `Super+Shift+Q`)
2. Press **`Super+Space`** — the semantic search bar appears
3. Type a query like `config` or `scripts` or `TODO`
4. Select a result — it opens in Neovim, Yazi, or the appropriate application

## Test It Out

Try these searches to see Ash in action:

| Query | What You'll Find |
|-------|------------------|
| `config` | Configuration files across your system |
| `scripts` | Shell scripts and automation files |
| `TODO` | Files containing task lists and notes |
| `files from 2h` | Everything modified in the last 2 hours |
| `files from 3d` | Everything modified in the last 3 days |

## Verify Installation

Run these commands to check that everything is working:

```bash
# Check Qdrant (vector database)
curl http://localhost:6333/health
# Expected: {"status":"ok","version":"..."}

# Check Ollama (AI model server)
curl http://localhost:11434/api/tags
# Expected: JSON listing "nomic-embed-text"

# Check LSFS daemon (file indexer)
systemctl --user is-active lsfs-daemon
# Expected: "active"

# Check launcher hook exists
test -x ~/.config/scripts/lsfs_launcher_hook.sh && echo "OK" || echo "MISSING"
```

## Something Not Working?

See the [Troubleshooting Guide](troubleshooting.md), or re-run the installer — it fixes most issues automatically:

```bash
curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash
```

---

**Next:** [How Search Works →](how-search-works.md) | [VMware Setup →](vmware-setup.md)
]]>

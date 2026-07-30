<![CDATA[# Backup & Restore

How to back up your Ash Linux data and recover from failures.

## What Needs Backing Up

Ash stores data in three places:

| Data | Location | What It Contains |
|------|----------|------------------|
| **Vector database** | `/var/lib/qdrant/` | All your file embeddings (searchable vectors) |
| **AI models** | `~/.ollama/` | Downloaded models (nomic-embed-text) |
| **Configuration** | `~/.config/scripts/`, `~/.config/hypr/` | Scripts, Hyprland config, daemon settings |
| **Service files** | `~/.config/systemd/user/` | systemd unit files for the daemon |

## Back Up

### Quick Backup (All Data)

```bash
# Create a timestamped backup of everything
mkdir -p ~/backups

# Vector database (your indexed file embeddings)
sudo tar -czf ~/backups/qdrant-$(date +%F).tar.gz /var/lib/qdrant

# AI models
tar -czf ~/backups/ollama-$(date +%F).tar.gz ~/.ollama

# Configuration and scripts
tar -czf ~/backups/ash-config-$(date +%F).tar.gz \
  ~/.config/scripts \
  ~/.config/systemd/user/lsfs-daemon.service \
  ~/.config/hypr \
  ~/.lsfsignore
```

### What You Can Skip

- **AI models** (`~/.ollama/`) — These can be re-downloaded with `ollama pull nomic-embed-text`
- **Vector database** (`/var/lib/qdrant/`) — The daemon will re-index your files automatically after a fresh install. This just saves time on large filesystems.

## Restore

### From Scratch (Disaster Recovery)

1. Start with a fresh Arch Linux VM
2. Run the Ash installer:
   ```bash
   curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash
   ```
3. Restore your backups:
   ```bash
   # Vector database
   sudo tar -xzf ~/backups/qdrant-*.tar.gz -C /

   # AI models (saves re-download time)
   tar -xzf ~/backups/ollama-*.tar.gz -C ~

   # Configuration
   tar -xzf ~/backups/ash-config-*.tar.gz -C ~
   ```
4. Restart services:
   ```bash
   sudo systemctl restart qdrant
   sudo systemctl restart ollama
   systemctl --user restart lsfs-daemon
   ```

### Restore Just the Config

If your services are running but your configuration was lost:

```bash
tar -xzf ~/backups/ash-config-*.tar.gz -C ~
systemctl --user restart lsfs-daemon
```

### Re-Index Without Backup

If you lost your vector database but still have your files, the daemon will re-index automatically:

```bash
# Restart the daemon — it will re-scan and re-index all watched paths
systemctl --user restart lsfs-daemon

# Watch it work:
journalctl --user -u lsfs-daemon -f
```

This takes a few minutes depending on how many files you have.

## What Survives a Reboot

Everything. All data is stored on the regular filesystem and persists across reboots. There is no ephemeral storage.

| Component | Persists? | Location |
|-----------|-----------|----------|
| Qdrant data | ✅ Yes | `/var/lib/qdrant/` |
| Ollama models | ✅ Yes | `~/.ollama/` |
| Launcher hook | ✅ Yes | `~/.config/scripts/lsfs_launcher_hook.sh` |
| LSFS daemon | ✅ Yes | `~/.config/scripts/lsfs_daemon.py` |
| Hyprland config | ✅ Yes | `~/.config/hypr/` |
| systemd services | ✅ Yes | `/etc/systemd/system/`, `~/.config/systemd/user/` |

---

**Next:** [Updating →](updating.md) | [Troubleshooting →](troubleshooting.md)
]]>

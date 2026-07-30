<![CDATA[# Updating Ash Linux

How to keep Ash and its components up to date.

## Update Everything

The easiest way to update Ash is to re-run the installer. It's idempotent — it updates all components without destroying your existing data:

```bash
curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash
```

This updates:
- ✅ Qdrant (downloads latest binary from GitHub releases)
- ✅ LSFS daemon and launcher hook (re-writes scripts)
- ✅ systemd service configurations
- ✅ Hyprland keybinding

This preserves:
- ✅ Your Qdrant vector database (indexed files)
- ✅ Your Ollama models
- ✅ Your personal configuration changes

## Update Individual Components

### Ollama

```bash
# Update Ollama itself
sudo pacman -S ollama

# Update the embedding model
ollama pull nomic-embed-text

# Verify
ollama list
```

### Qdrant

The installer automatically fetches the latest Qdrant release. To update Qdrant independently:

```bash
# Re-run the installer (safest method)
curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash
```

### LSFS Daemon & Launcher Hook

These are updated automatically when you re-run the installer. The scripts at `~/.config/scripts/` are overwritten with the latest versions.

### Arch System Packages

Keep your base system updated with the usual Arch commands:

```bash
sudo pacman -Syu
```

If system updates break any Ash services, re-run the installer to fix them:

```bash
curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash
```

## Checking Current Versions

```bash
# Qdrant version
curl -s http://localhost:6333 | grep version

# Ollama version
curl -s http://localhost:11434/api/version

# LSFS daemon status
systemctl --user status lsfs-daemon
```

---

**Next:** [Backup & Restore →](backup-and-restore.md) | [Troubleshooting →](troubleshooting.md)
]]>

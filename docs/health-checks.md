<![CDATA[# Health Checks

Quick commands to verify every Ash component is running correctly.

## All-in-One Check

Copy and paste this to check everything at once:

```bash
echo "=== Ash Linux Health Check ==="
echo ""
echo -n "Qdrant:        "; curl -sf http://localhost:6333/health > /dev/null 2>&1 && echo "✅ Running" || echo "❌ Down"
echo -n "Ollama:        "; curl -sf http://localhost:11434/api/tags > /dev/null 2>&1 && echo "✅ Running" || echo "❌ Down"
echo -n "LSFS Daemon:   "; systemctl --user is-active lsfs-daemon > /dev/null 2>&1 && echo "✅ Running" || echo "❌ Down"
echo -n "Launcher Hook: "; test -x ~/.config/scripts/lsfs_launcher_hook.sh && echo "✅ Present" || echo "❌ Missing"
echo -n "Embed Model:   "; curl -sf http://localhost:11434/api/tags 2>/dev/null | grep -q nomic-embed-text && echo "✅ Loaded" || echo "⚠️  Not found"
echo -n "Qdrant Data:   "; curl -sf http://localhost:6333/collections/apps > /dev/null 2>&1 && echo "✅ Collection exists" || echo "⚠️  No collection"
```

## Individual Checks

### Qdrant (Vector Database)

```bash
# Health
curl http://localhost:6333/health
# Expected: {"status":"ok","version":"..."}

# Collection info (shows how many files are indexed)
curl http://localhost:6333/collections/apps
```

### Ollama (AI Model Server)

```bash
# Version
curl http://localhost:11434/api/version

# Check the embedding model is available
curl http://localhost:11434/api/tags | grep nomic-embed-text

# If missing, pull it:
ollama pull nomic-embed-text
```

### LSFS Daemon (File Indexer)

```bash
# Status
systemctl --user status lsfs-daemon

# Live logs (watch it index files in real-time)
journalctl --user -u lsfs-daemon -f
```

### Launcher Hook

```bash
# Check it exists and is executable
ls -l ~/.config/scripts/lsfs_launcher_hook.sh

# Test it manually (opens a search prompt)
bash ~/.config/scripts/lsfs_launcher_hook.sh
```

## Fix Everything

If anything is broken, re-run the installer — it's idempotent and will repair all services:

```bash
curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash
```

---

**Next:** [Troubleshooting →](troubleshooting.md) | [Quick Start →](quickstart.md)
]]>

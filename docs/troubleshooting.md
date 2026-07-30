<![CDATA[# Troubleshooting

Solutions for common issues with Ash Linux.

## Quick Fix: Re-Run the Installer

Most issues can be fixed by re-running the Ash installer. It's designed to be idempotent and will repair broken services automatically:

```bash
curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash
```

---

## Search Issues

### Super+Space does nothing

**Cause:** The keybinding isn't configured in Hyprland, or the launcher hook script is missing.

**Fix:**
```bash
# Check if the launcher hook exists and is executable
ls -l ~/.config/scripts/lsfs_launcher_hook.sh

# Check if Hyprland has the keybinding
grep "Super.*Space" ~/.config/hypr/hyprland.conf

# If either is missing, re-run the installer
curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash
```

### Search returns no results

**Cause:** Either the services aren't running, or no files have been indexed yet.

**Fix:**
```bash
# Check if Qdrant is running
curl http://localhost:6333/health

# Check if Ollama is running
curl http://localhost:11434/api/tags

# Check if the LSFS daemon is indexing
systemctl --user status lsfs-daemon

# If services are down, start them:
sudo systemctl start qdrant
sudo systemctl start ollama
systemctl --user start lsfs-daemon
```

### "Ollama not running" notification

**Fix:**
```bash
sudo systemctl start ollama
sudo systemctl enable ollama

# Verify
curl http://localhost:11434/api/tags | grep nomic-embed-text

# If model is missing, pull it:
ollama pull nomic-embed-text
```

### "Qdrant not running" notification

**Fix:**
```bash
sudo systemctl start qdrant
sudo systemctl enable qdrant

# Verify
curl http://localhost:6333/health
```

If Qdrant fails to start:
```bash
# Check logs
sudo journalctl -u qdrant -n 30

# Common issue: port already in use
sudo lsof -i :6333

# Common issue: storage permissions
sudo ls -la /var/lib/qdrant/
```

---

## Desktop Issues

### Hyprland won't start (black screen)

**Most common cause:** Missing VMX display fix (VMware only).

**Fix:**
1. Shut down the VM
2. Add to the `.vmx` file on your **host** machine:
   ```
   mks.enableVulkanRenderer = "FALSE"
   svga.disableFIFO = "TRUE"
   ```
3. Boot the VM

If not on VMware, check the Hyprland log:
```bash
# Switch to TTY2 with Ctrl+Alt+F2
cat ~/.local/share/hyprland/hyprland.log | tail -50
```

### Clipboard not working (VMware)

**Fix:** Add to the `.vmx` file on your **host** machine:
```
isolation.tools.copy.disable = "FALSE"
isolation.tools.paste.disable = "FALSE"
isolation.tools.setGUIOptions.enable = "TRUE"
```

Then restart the guest tools inside the VM:
```bash
sudo systemctl restart open-vm-tools
```

### Screen tearing

Add to the `.vmx` file:
```
mks.enableVulkanRenderer = "FALSE"
svga.disableFIFO = "TRUE"
```

---

## Service Issues

### LSFS daemon not running

```bash
# Check status
systemctl --user status lsfs-daemon

# Check logs for errors
journalctl --user -u lsfs-daemon -n 50

# Restart
systemctl --user restart lsfs-daemon

# If it still fails, check the daemon script exists
ls -l ~/.config/scripts/lsfs_daemon.py
```

### LSFS daemon not indexing new files

```bash
# Check what the daemon is doing
journalctl --user -u lsfs-daemon -f

# Check the Qdrant collection
curl -s http://localhost:6333/collections/apps | python3 -m json.tool

# Verify inotify limits are high enough
cat /proc/sys/fs/inotify/max_user_watches
# Should be at least 524288
```

### Services don't start on boot

```bash
# Enable all services
sudo systemctl enable qdrant
sudo systemctl enable ollama
sudo systemctl enable open-vm-tools
systemctl --user enable lsfs-daemon

# Enable lingering (for user services to start before login)
loginctl enable-linger $(whoami)
```

---

## Health Check (All-in-One)

Run this to check the status of every component:

```bash
echo "=== Ash Health Check ==="
echo ""

echo -n "Qdrant:      "
curl -sf http://localhost:6333/health > /dev/null 2>&1 && echo "✅ OK" || echo "❌ DOWN"

echo -n "Ollama:      "
curl -sf http://localhost:11434/api/tags > /dev/null 2>&1 && echo "✅ OK" || echo "❌ DOWN"

echo -n "LSFS Daemon: "
systemctl --user is-active lsfs-daemon > /dev/null 2>&1 && echo "✅ OK" || echo "❌ DOWN"

echo -n "Launcher:    "
test -x ~/.config/scripts/lsfs_launcher_hook.sh && echo "✅ OK" || echo "❌ MISSING"

echo -n "Model:       "
curl -sf http://localhost:11434/api/tags 2>/dev/null | grep -q nomic-embed-text && echo "✅ nomic-embed-text loaded" || echo "⚠️  Model not found"

echo ""
echo "To fix issues, re-run: curl -sfL https://raw.githubusercontent.com/exonew2/files/main/scripts/ultimate-fix-v2.sh | sudo bash"
```

---

## Still Stuck?

1. **Re-run the installer** — it fixes 90% of issues automatically
2. **Check the logs** — `journalctl --user -u lsfs-daemon -f` and `sudo journalctl -u qdrant -f`
3. **Open an issue** — [GitHub Issues](https://github.com/exonew2/files/issues)

---

**Next:** [Health Checks →](health-checks.md) | [Quick Start →](quickstart.md)
]]>

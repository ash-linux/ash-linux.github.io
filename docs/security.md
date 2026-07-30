<![CDATA[# Security

How Ash Linux protects your data and what to be aware of.

## Security Model

Ash is designed around **local-only operation**. No data ever leaves your machine.

### Network Isolation

All services listen on **localhost only** (`127.0.0.1`):

| Service | Bind Address | Port |
|---------|-------------|------|
| Ollama | `127.0.0.1` | `11434` |
| Qdrant | `127.0.0.1` | `6333` |

No ports are exposed to the network. External machines cannot reach these services.

### Process Isolation

Services are sandboxed via systemd:

| Service | Sandboxing |
|---------|-----------|
| **Qdrant** | Runs as dedicated `qdrant` user, CPUQuota=50%, MemoryMax=2G |
| **LSFS Daemon** | Runs as your user, Nice=19 (lowest CPU priority), IOSchedulingClass=idle |
| **Launcher Hook** | Runs as your user, no elevated privileges |

### VM Boundary

When running in a VM (recommended), the VM boundary provides additional isolation from your host machine. Ash services inside the VM cannot access host files, network, or devices without explicit configuration.

## Known Gaps

### No Firewall

Ash does not configure a firewall inside the VM by default. While all services bind to localhost, adding a firewall provides defense-in-depth:

```bash
# Optional: enable the included firewall
sudo systemctl enable --now firewalld
sudo firewall-cmd --zone=public --remove-service=ssh  # if you don't need SSH
```

### No AppArmor/SELinux

Ash does not ship with mandatory access control profiles. Arch Linux uses neither AppArmor nor SELinux by default.

### Clipboard Sharing

Enabling clipboard sharing between host and VM requires disabling VMware's clipboard isolation in the `.vmx` file. This is a deliberate trade-off between convenience and isolation. If you don't need clipboard sharing, leave the default isolation in place.

## Best Practices

1. **Run Ash in a VM** — provides hardware-level isolation from your host
2. **Don't expose ports** — keep all services on localhost
3. **Back up your data** — see [Backup & Restore](backup-and-restore.md)
4. **Keep the system updated** — `sudo pacman -Syu` regularly
5. **Don't run untrusted models** — only use embedding models from trusted sources (Ollama's official library)

## Reporting Security Issues

If you find a security vulnerability, please report it privately via [GitHub Security Advisories](https://github.com/exonew2/files/security/advisories) rather than opening a public issue.

---

**Next:** [Comparison →](comparison.md) | [Troubleshooting →](troubleshooting.md)
]]>

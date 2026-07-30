<![CDATA[# Security Policy

## Security Model

Ash Linux is designed for **local-only operation**. All services bind to `127.0.0.1` and no data is transmitted over the network.

- **Ollama**: `127.0.0.1:11434` — local AI model server
- **Qdrant**: `127.0.0.1:6333` — local vector database
- **LSFS Daemon**: user-space process, no network ports
- **Launcher Hook**: user-space bash script, no elevated privileges

### Sandboxing

- `qdrant.service`: Runs as dedicated `qdrant` user with CPUQuota=50%, MemoryMax=2G
- `lsfs-daemon.service`: Runs at lowest CPU/IO priority (Nice=19, IOSchedulingClass=idle)
- All services restart on failure via systemd

### VM Isolation

When running in a VM (recommended), the hypervisor boundary provides additional isolation from the host system.

## Known Gaps

- No firewall configured inside the VM by default (`firewalld` is installed but not enabled)
- No AppArmor or SELinux profiles
- Enabling clipboard sharing requires disabling VMware's clipboard isolation

## Reporting Vulnerabilities

**Please do not open public issues for security vulnerabilities.**

Instead, report them privately via:
- [GitHub Security Advisories](https://github.com/exonew2/files/security/advisories)
- Email: security@ash.sh (if available)

We will acknowledge receipt within 48 hours and provide a fix timeline.

## Supported Versions

| Version | Supported |
|---------|-----------|
| v2.x (current) | ✅ Active |
| v1.x | ❌ No longer supported |
]]>

# Security Policy

## Reporting a Vulnerability

**Do not file public issues for security vulnerabilities.**

Email: **security@ash.sh**

Include:
- Description of the vulnerability
- Steps to reproduce
- Affected versions
- Potential impact
- Suggested fix (if any)

We respond within 48 hours. Critical vulnerabilities get a patch release within 7 days.

## Scope

In scope:
- ISO build process (supply chain)
- First-boot scripts (privilege escalation)
- Guest agent integrations (VM escape)
- Firewall rules (network isolation)
- Snapshot/rollback mechanism (data integrity)
- Signature verification (supply chain)

Out of scope:
- Upstream Arch Linux packages (report to Arch)
- Upstream AI tools (Ollama, llama.cpp, Qdrant — report to their projects)
- Hypervisor vulnerabilities (report to VMware/Oracle/Parallels/QEMU)
- Hardware vulnerabilities (Spectre, Meltdown, etc.)

## Supply Chain Security

Every release provides:
1. **SHA256** — Basic integrity
2. **minisign** — Ed25519 signatures, fast offline verification
3. **cosign (keyless)** — Sigstore transparency log, OIDC identity
4. **SLSA Level 3 Provenance** — Build attestation via GitHub Actions

Verify **all three** before booting.

## Threat Model

| Threat | Mitigation |
|--------|------------|
| Malicious ISO on mirror | minisign + cosign + SLSA verification |
| Compromised build runner | SLSA L3 provenance, reproducible builds |
| VM escape via guest agent | Minimal agents, seccomp, no host fs access |
| AI model supply chain | Models pulled from Ollama library only, user consent |
| Persistent compromise | Btrfs snapshots + rollback, Qdrant excluded from snaps |

## Security Updates

- Weekly auto-update timer (opt-in) creates pre-snapshot
- Kernel held (`IgnorePkg = linux linux-headers linux-firmware`)
- Security-only updates via `pacman -Syu --ignore=linux*`
- Emergency releases for critical CVEs within 24h

## Responsible Disclosure Timeline

| Severity | Response | Patch Release |
|----------|----------|---------------|
| Critical (CVSS 9-10) | 24 hours | 7 days |
| High (CVSS 7-8.9) | 48 hours | 14 days |
| Medium (CVSS 4-6.9) | 1 week | 30 days |
| Low (CVSS 0-3.9) | 2 weeks | Next scheduled |

## Hall of Fame

Security researchers who responsibly disclosed:
- *(awaiting first submission)*
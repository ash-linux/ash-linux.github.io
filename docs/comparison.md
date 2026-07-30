<![CDATA[# Comparison — Ash Linux vs Alternatives

How Ash compares to other approaches for setting up an AI-powered Linux desktop.

## At a Glance

| Feature | **Ash Linux** | Plain Arch Setup | Ubuntu Desktop | ISO-Based AI OS |
|---------|---------------|------------------|----------------|-----------------|
| **Setup Time** | ~2 minutes (one-liner) | 2–4 hours | 30 min + manual config | Download ISO + boot |
| **Semantic Search** | ✅ Built-in | ❌ None | ❌ None | ⚠️ Varies |
| **Vector Database** | ✅ Qdrant (auto-configured) | ❌ None | ❌ None | ⚠️ Varies |
| **AI Launcher** | ✅ Super+Space | ❌ None | ❌ None | ❌ Usually absent |
| **One-Line Install** | ✅ `curl \| sudo bash` | ❌ N/A | ❌ N/A | ❌ ISO download |
| **100% Offline** | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |
| **VMware Optimized** | ✅ Documented workarounds | ❌ Manual | ❌ Manual | ❌ Generic |
| **Customizable** | ✅ Full Arch control | ✅ Full Arch control | ❌ Constrained | ⚠️ Fixed filesystem |

## Ash vs Plain Arch Linux

If you already run Arch, Ash adds a complete AI search stack in one command:

| | Ash Linux | Doing It Yourself |
|--|-----------|-------------------|
| Time to set up | ~2 minutes | 2–4 hours of reading docs, installing, configuring |
| Semantic file search | Pre-configured and working | Doesn't exist — you'd have to build it |
| Vector database | Qdrant systemd service, auto-starts | Manual download, configuration, service setup |
| AI launcher | `Super+Space` → search by meaning | Nothing comparable |
| Auto-login | Configured out of the box | Manual agetty + `.bash_profile` setup |
| Hyprland for VMware | VMX workarounds documented | Trial and error |

## Ash vs Ubuntu Desktop

| | Ash Linux | Ubuntu |
|--|-----------|--------|
| Package freshness | Arch (rolling release) | Ubuntu (fixed release every 6 months) |
| Semantic search | Built-in, native | None (GNOME Search is filename-only) |
| AI stack | One-command deploy | Manual pip/apt install + configuration |
| Wayland compositor | Hyprland (tiling, pre-configured) | GNOME (generic, stacking) |
| Storage overhead | Minimal (scripts only) | 5–10 GB base system |
| Learning curve | Requires Arch familiarity | More beginner-friendly |

## Ash vs ISO-Based AI OS Projects

| | Ash Linux | ISO-Based Projects |
|--|-----------|-------------------|
| Deployment | Script on existing Arch | Download + boot ISO |
| Destructive? | No — adds to your existing setup | Yes — new VM or partition |
| Filesystem | Your existing setup, untouched | Btrfs subvolumes, read-only root |
| Updates | Re-run the installer script | ISO re-download or A/B updates |
| Snapshots | Manual backups | Built-in (Snapper, Btrfs) |
| Flexibility | Full control of base system | Opinionated, pre-baked defaults |

## When to Use What

| Your Situation | Recommendation |
|----------------|----------------|
| Already on Arch + Hyprland, want semantic search | **Ash Linux** (one-liner install) |
| Setting up a fresh VM for AI development | Fresh Arch install → Ash one-liner |
| Want a zero-config AI appliance | Consider an ISO-based AI OS project |
| Want maximum customization | Plain Arch + cherry-pick Ash scripts |
| Ephemeral experiments | Ash in a scratchable VM |

## What Makes Ash Different

- **Semantic search via bash** — The launcher is ~50 lines of pure bash using `curl`. No Python in the search path.
- **nomic-embed-text** — 768-dim, Ollama-native embedding model. Small, fast, runs on CPU.
- **Qdrant standalone binary** — No AUR, no Docker, no Python SDK. Static binary from GitHub.
- **VMware-first** — Display and clipboard workarounds are documented and automated.
- **Idempotent installer** — Re-run anytime. It fixes, never breaks.

## Limitations

Be aware of these trade-offs:

- **Clipboard** needs a manual host-side `.vmx` edit for VMware
- **Display** needs a manual host-side `.vmx` edit for Hyprland stability
- **Default index scope** is limited — configure additional paths manually
- **No automatic rollback** — no Btrfs snapshots; back up manually
- **Arch-only** — the installer requires pacman; no Ubuntu/Fedora support yet

---

**Next:** [Quick Start →](quickstart.md) | [How Search Works →](how-search-works.md)
]]>

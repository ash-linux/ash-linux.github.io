<![CDATA[# GPU & Performance

Getting the most out of Ash Linux — GPU passthrough, memory tuning, and performance tips.

## GPU Support

### Does Ash Need a GPU?

**No.** The `nomic-embed-text` embedding model runs comfortably on CPU. Semantic search works at full speed without any GPU.

A GPU provides modest speedup when batch-indexing many files at once, but for normal usage (single queries, real-time indexing), CPU is more than sufficient.

### GPU Acceleration in VMware

If you want GPU acceleration inside a VMware VM:

1. Enable **3D acceleration** in VM Settings → Display
2. Allocate at least **4 GB video memory**
3. Set at least **4 CPU cores**

> **Note:** You still need the VMX display workaround for Hyprland stability. See [VMware Setup](vmware-setup.md).

### Ollama GPU Detection

Ollama automatically detects and uses available GPU acceleration:

- **NVIDIA** → CUDA (install `nvidia` drivers inside the VM)
- **AMD** → ROCm
- **Apple Silicon** → Metal (macOS host only)

No manual configuration needed. Ollama picks the best available backend automatically.

### Verify GPU Usage

```bash
# Check Ollama version and backend
curl http://localhost:11434/api/version

# Benchmark embedding speed
time curl -X POST http://localhost:11434/api/embeddings \
  -d '{"model":"nomic-embed-text","prompt":"benchmark test"}'

# Check NVIDIA GPU (if available)
nvidia-smi

# Check AMD GPU (if available)
rocm-smi
```

## Performance Tuning

### Memory Allocation

| Component | Default Limit | Recommendation |
|-----------|---------------|----------------|
| **Qdrant** | 2 GB (systemd MemoryMax) | Sufficient for up to ~1M files |
| **Ollama** | ~1 GB (model size) | Increases with larger models |
| **LSFS Daemon** | No hard limit (low priority) | Uses minimal memory |

Total system recommendation: **8 GB RAM** for comfortable operation.

### CPU Allocation

| Component | Priority | Notes |
|-----------|----------|-------|
| **Qdrant** | CPUQuota=50% | Prevents Qdrant from starving other processes |
| **LSFS Daemon** | Nice=19, IOSchedulingClass=idle | Runs at lowest priority |
| **Ollama** | Normal | Handles embedding requests on demand |

### Filesystem Tuning

The installer automatically sets these sysctl values:

```bash
# Maximum inotify watches (for the LSFS daemon)
fs.inotify.max_user_watches = 524288

# Lower swappiness (keep more in RAM)
vm.swappiness = 10
```

### Indexing Speed

The LSFS daemon processes files one at a time with debouncing. For large initial imports:

1. **Let it run** — the daemon will work through the backlog at its own pace
2. **Monitor progress** — `journalctl --user -u lsfs-daemon -f`
3. **Don't restart during indexing** — let it finish the current batch

Typical indexing speed: ~5-10 files/second on CPU, ~20-50 files/second with GPU.

---

**Next:** [How Search Works →](how-search-works.md) | [Cluster Mode →](cluster.md)
]]>

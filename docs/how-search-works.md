<![CDATA[# How Search Works

A plain-language explanation of what happens when you press `Super+Space` and search for a file in Ash Linux.

## The Short Version

1. You type a query in natural language
2. Ash converts your words into a mathematical representation of their *meaning*
3. Ash compares that meaning against every file it has indexed
4. The most similar files appear as results

No keyword matching. No regex. Just "find files that are *about* this."

## Step by Step

### 1. You Press `Super+Space`

A clean search bar (powered by **Wofi**) appears on your screen. You type a query — anything from a single word like `config` to a full sentence like `the backup script I wrote last week`.

### 2. Your Query Becomes a Vector

Your text is sent to **Ollama**, which runs a model called `nomic-embed-text` on your local machine. This model converts your query into a **768-dimensional vector** — essentially a list of 768 numbers that represent the *meaning* of what you typed.

Think of it like GPS coordinates, but for meaning. The word "backup" would be close to "snapshot" and "archive" in this space, but far from "cooking recipe."

```
"find my config files" → [0.023, -0.156, 0.891, ..., 0.042]  (768 numbers)
```

This happens entirely on your machine. Nothing is sent to the internet.

### 3. Qdrant Finds Similar Files

That vector is sent to **Qdrant**, a vector database running on your machine. Qdrant already has vectors for every file on your system (the LSFS daemon creates these in the background — more on that below).

Qdrant uses **cosine similarity** to find the files whose vectors are most similar to your query vector. This is mathematically identical to "find files whose content is closest in *meaning* to what the user typed."

### 4. Results Appear

The top matches are displayed in a selection menu. Each result shows the file path and a similarity score. Select any result to open it:

- **Text/code files** → open in **Neovim** (inside Kitty terminal)
- **Directories** → open in **Yazi** file manager
- **`.desktop` files** → launch the application

### 5. Time-Based Fallback

If your query looks like a time expression (e.g., `files from 2h`, `past 3 days`), Ash skips vector search entirely and uses `fd` or `find` to locate recently modified files. This works even if Qdrant is down.

## How Files Get Indexed

The **LSFS daemon** runs silently in the background as a systemd user service. It:

1. **Watches your filesystem** using Linux's `inotify` system
2. **Detects changes** — new files, edits, moves, deletes
3. **Reads file content** and sends it to Ollama for embedding
4. **Stores the vectors** in Qdrant's `apps` collection

This means every file you create or edit becomes searchable within seconds — automatically, silently, locally.

### What Gets Indexed

By default, the daemon indexes `~/.config/scripts` and other configured paths. You can customize this.

### What Gets Ignored

The `~/.lsfsignore` file controls what the daemon skips (uses the same syntax as `.gitignore`):

```
# Already ignored by default:
node_modules/
__pycache__/
.git/
*.pyc
*.mp3 *.mp4 *.png *.jpg     # binary/media files
*.zip *.tar *.gz             # archives
*.o *.so *.dll *.exe         # compiled binaries
.venv/ venv/ target/ build/  # build artifacts
```

Edit `~/.lsfsignore` to add or remove ignore rules.

## The Components

| Component | Role | Runs As |
|-----------|------|---------|
| **Wofi** | Search prompt UI | Triggered by keyboard shortcut |
| **Ollama** | Generates text embeddings | `systemd` system service on port `11434` |
| **nomic-embed-text** | The AI model that understands text | Loaded by Ollama, pinned in memory |
| **Qdrant** | Stores and searches file vectors | `systemd` system service on port `6333` |
| **LSFS Daemon** | Watches filesystem, indexes changes | `systemd` user service |
| **Launcher Hook** | Connects everything together | Pure bash script (~50 lines) |

## Key Design Decisions

### Why is the launcher a bash script?

The `Super+Space` search pipeline is a pure-bash script of ~50 lines. Its only dependency is `curl`. This means:
- No Python in the search path = no import delays, no virtual environments
- Instant startup, zero overhead
- Easy to audit, modify, or debug

### Why nomic-embed-text?

- **Ollama-native** — no separate Python ML stack needed, no PyTorch
- **768 dimensions** — good balance of accuracy and speed
- **No prefixes** — other models require `search_query:` or `passage:` prefixes; nomic works with raw text
- **CPU-friendly** — runs smoothly without a GPU

### Why Qdrant as a standalone binary?

- **No Docker dependency** — runs as a simple binary
- **No AUR builds** — downloaded directly from GitHub releases
- **Simple systemd service** — starts at boot, auto-restarts on failure

---

**Next:** [VMware Setup →](vmware-setup.md) | [Backup & Restore →](backup-and-restore.md)
]]>

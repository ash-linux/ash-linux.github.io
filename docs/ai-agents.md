<![CDATA[# AI Agent Integration

Connect your AI coding assistants to Ash's knowledge graph and memory system.

## Supported Agents

Ash includes a one-command configurator that sets up the following AI tools:

| Agent | Configuration | What It Gets |
|-------|--------------|--------------|
| **Claude Desktop** | MCP server in `claude_desktop_config.json` | Codebase memory via MCP protocol |
| **Cursor** | `.cursorrules` file | Project context rules |
| **Windsurf** | `.windsurfrules` file | Project context rules |
| **Cline** (VS Code) | `.vscode/cline_mcp.json` | MCP server + DB endpoint |
| **Gemini** | `gemini-context.json` | System instructions + endpoints |

## One-Command Setup

```bash
# Configure all agents at once
bash ai-services/configure-agents.sh
```

This script:
1. Detects which agents are available on your system
2. Creates or merges configuration files for each
3. Connects them to the shared codebase memory database
4. Reports what was configured, skipped, or failed

## What the Agents Get Access To

### Codebase Memory (MCP)

A knowledge graph that provides persistent context across AI sessions:

- **Semantic search** over your codebase
- **Entity resolution** — understands relationships between files, functions, and concepts
- **Memory persistence** — AI agents remember context from previous sessions

Backed by a local SQLite database and vector embeddings.

### Database Endpoint

Direct access to the prompt storage and graph persistence layer:

- **Path:** `ai-services/data/memory.db` (SQLite)
- **Vectors:** `ai-services/data/vectors/` (ChromaDB)

### Project Context Rules

For agents that use rules files (Cursor, Windsurf), the configurator injects:

- Database paths for persistent context
- Vector storage locations
- Project-specific coding guidelines

## Manual Configuration

### Claude Desktop

Add to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `~/.config/Claude/claude_desktop_config.json` (Linux):

```json
{
  "mcpServers": {
    "codebase-memory": {
      "command": "codebase-memory-mcp",
      "args": [],
      "env": {
        "CBM_MEMORY_PATH": "/path/to/ai-services/data/memory.db"
      }
    }
  }
}
```

### Cursor

The configurator creates/updates `.cursorrules` in the project root with memory system paths and project context.

### Cline (VS Code)

The configurator creates `.vscode/cline_mcp.json` with MCP server and database endpoint configurations.

## Requirements

- `codebase-memory-mcp` must be installed (`npm install -g codebase-memory-mcp` or equivalent)
- Python 3 (for the DB endpoint and configurator)
- The respective AI tools installed on your system

---

**Next:** [Security →](security.md) | [How Search Works →](how-search-works.md)
]]>

<![CDATA[# Contributing to Ash Linux

Thanks for your interest in contributing! Here's how to get involved.

## Ways to Contribute

### 🐛 Report Bugs
Found something broken? [Open an issue](https://github.com/exonew2/files/issues) with:
- What you expected to happen
- What actually happened
- Steps to reproduce
- Output of the health check (see [Health Checks](docs/health-checks.md))

### 💡 Suggest Features
Have an idea? [Open a discussion](https://github.com/exonew2/files/discussions) or issue describing:
- What problem it solves
- How you envision it working
- Whether you'd be willing to help implement it

### 🔧 Submit Code
1. Fork the repo
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Run syntax checks:
   ```bash
   # Shell scripts
   bash -n scripts/*.sh

   # Python scripts
   python -m py_compile ai-services/*.py
   ```
5. Submit a pull request

## Areas Where Help Is Most Needed

| Area | Description |
|------|-------------|
| **VMware clipboard automation** | Host-side scripting to auto-configure `.vmx` clipboard settings |
| **Additional embedding models** | Support for switching between different Ollama embedding models |
| **Qdrant backup/restore** | Automated backup and restore tooling for the vector database |
| **Ubuntu/Fedora support** | Porting the installer to work with `apt` and `dnf` |
| **Broader index paths** | Configurable watched directories with a settings UI |
| **Image/PDF indexing** | Extracting text from images and PDFs for vector search |

## Code Style

- **Shell scripts**: Use `set -euo pipefail`, consistent logging functions
- **Python**: Follow PEP 8, use type hints where practical
- **Documentation**: Write for end users, not developers. Explain *why*, not just *how*.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
]]>

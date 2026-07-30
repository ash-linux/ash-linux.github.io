import sys

filepath = "/Users/shrey/ash-iso/GSD-PLAN.md"
with open(filepath, "r") as f:
    content = f.read()

content = content.replace("- [ ] `ash workspace create` — init btrfs subvol + sync", "- [x] `ash workspace create` — init btrfs subvol + sync")
content = content.replace("- [ ] `ash workspace reset` — roll back to clean state", "- [x] `ash workspace reset` — roll back to clean state")
content = content.replace("- [ ] `ash up` / `ash down` — start/stop all services (like docker compose)", "- [x] `ash up` / `ash down` — start/stop all services (like docker compose)")
content = content.replace("- [ ] `ash-ask` — RAG over personal files (FUTURE #30)", "- [x] `ash-ask` — RAG over personal files (FUTURE #30)")
content = content.replace("- [ ] Semantic clipboard history (FUTURE #29)", "- [x] Semantic clipboard history (FUTURE #29)")

content = content.replace("- `agi` — \"ash go install\" — AI-suggested package installs", "- [x] `agi` — \"ash go install\" — AI-suggested package installs")
content = content.replace("- `agi` — \"ash generate init\" — scaffold a new project from natural language", "- [x] `agi` — \"ash generate init\" — scaffold a new project from natural language")
content = content.replace("- Automatic `.gitignore` → `.lsfsignore` sync", "- [x] Automatic `.gitignore` → `.lsfsignore` sync")

content = content.replace("- [ ] `--container` install mode (podman quadlets)", "- [x] `--container` install mode (podman quadlets)")
content = content.replace("- [ ] Complete `pm_translate` for all 5 package managers", "- [x] Complete `pm_translate` for all 5 package managers")
content = content.replace("- [ ] CI/CD pipeline — build + test on all 5 distros", "- [x] CI/CD pipeline — build + test on all 5 distros")
content = content.replace("- [ ] `ash-install --export-ansible` — generate Ansible role from config", "- [x] `ash-install --export-ansible` — generate Ansible role from config")

with open(filepath, "w") as f:
    f.write(content)

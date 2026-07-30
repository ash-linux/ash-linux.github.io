import sys
filepath = "/Users/shrey/ash-iso/FUTURE.md"
with open(filepath, "r") as f:
    lines = f.readlines()

updates = {
    "### 31. Semantic clipboard history": "*STATUS: ✅ Done*\n",
    "### 101. `ash-ask` — conversational file Q&A": "*STATUS: ✅ Done*\n",
    "### 57. Containerized deployment (--container)": "*STATUS: ✅ Done*\n",
    "### 64. Ansible role generation": "*STATUS: ✅ Done*\n",
    "### 96. `agi` — AI-assisted package installation": "*STATUS: ✅ Done*\n",
    "### 97. `ash workspace` — dev environments": "*STATUS: ✅ Done*\n"
}

# The wording of the headers might slightly differ. I'll search via string matching.
new_lines = []
for i, line in enumerate(lines):
    new_lines.append(line)
    if "### 31. Semantic clipboard history" in line:
        new_lines.append("*STATUS: ✅ Done*\n")
    if "### 101. `ash-ask`" in line:
        new_lines.append("*STATUS: ✅ Done*\n")
    if "### 57. Containerized deployment" in line:
        new_lines.append("*STATUS: ✅ Done*\n")
    if "### 64. Ansible role generation" in line:
        new_lines.append("*STATUS: ✅ Done*\n")
    if "### 96. `agi`" in line:
        new_lines.append("*STATUS: ✅ Done*\n")
    if "### 97. `ash workspace`" in line:
        new_lines.append("*STATUS: ✅ Done*\n")

with open(filepath, "w") as f:
    f.writelines(new_lines)

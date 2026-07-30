import sys, os

INSTALLER_PATH = "/Users/shrey/ash-iso/scripts/ash-install.sh"
with open(INSTALLER_PATH, "r") as f:
    lines = f.readlines()

new_lines = []
in_pm_translate = False
skip_pm_translate = False

for i, line in enumerate(lines):
    if "WEBHOOK_URL=\"\"; AUDIT_LEVEL=\"none\"; MAX_RETRIES=5" in line:
        new_lines.append(line)
        new_lines.append("CONTAINER_MODE=false; EXPORT_ANSIBLE=false\n")
        continue

    if "--help)" in line:
        new_lines.append("  --container) CONTAINER_MODE=true ;;\n")
        new_lines.append("  --export-ansible) EXPORT_ANSIBLE=true ;;\n")
        new_lines.append(line.replace("[--help]", "[--help] [--container] [--export-ansible]"))
        continue

    if "declare -A PKG_MAP=(" in line:
        in_pm_translate = True
        skip_pm_translate = True
        
        # Inject new pm_translate maps
        maps = """  declare -A PKG_MAP=(
    [python]=python3 [python-pip]=python3-pip [curl]=curl [wget]=wget [jq]=jq [fd]=fd-find
    [openssh]=openssh-server [kitty]=kitty [lm_sensors]=lm-sensors [python-requests]=python3-requests
    [python-aiohttp]=python3-aiohttp [swaync]=swaync [wofi]=wofi [fuse3]=fuse3 [inotify-tools]=inotify-tools
    [hyprland]=hyprland [waybar]=waybar [rofi]=rofi [go]=go [rust]=rust [node]=nodejs
    [ollama]=ollama [qdrant]=qdrant
  )
  declare -A DNF_MAP=(
    [python]=python3 [python-pip]=python3-pip [fd]=fd-find [openssh]=openssh-server
    [lm_sensors]=lm_sensors [python-requests]=python3-requests [swaync]=swaync
    [wofi]=wofi [fuse3]=fuse3 [kitty]=kitty [inotify-tools]=inotify-tools
    [hyprland]=hyprland [waybar]=waybar [rofi]=rofi [go]=golang [rust]=rust [node]=nodejs
    [ollama]=ollama [qdrant]=qdrant
  )
  declare -A APT_MAP=(
    [fd]=fd-find [openssh]=openssh-server [lm_sensors]=lm-sensors [python-requests]=python3-requests
    [python-aiohttp]=python3-aiohttp [swaync]=swaync [wofi]=wofi [fuse3]=fuse3 [inotify-tools]=inotify-tools
    [hyprland]=hyprland [waybar]=waybar [rofi]=rofi [go]=golang-go [rust]=rustc [node]=nodejs
    [ollama]=ollama [qdrant]=qdrant
  )
  declare -A ZYPPER_MAP=(
    [fd]=fd-find [openssh]=openssh-server [lm_sensors]=lm_sensors [python-requests]=python3-requests
    [python-aiohttp]=python3-aiohttp [swaync]=swaync [wofi]=wofi [inotify-tools]=inotify-tools
    [hyprland]=hyprland [waybar]=waybar [rofi]=rofi [go]=go [rust]=rust [node]=nodejs
    [ollama]=ollama [qdrant]=qdrant
  )
  declare -A APK_MAP=(
    [fd]=fd [openssh]=openssh [lm_sensors]=lm-sensors [python-requests]=py3-requests
    [python-aiohttp]=py3-aiohttp [swaync]=swaync [wofi]=wofi [inotify-tools]=inotify-tools
    [hyprland]=hyprland [waybar]=waybar [rofi]=rofi [go]=go [rust]=rust [node]=nodejs
    [ollama]=ollama [qdrant]=qdrant
  )
"""
        new_lines.append(maps)
        continue
    
    if skip_pm_translate:
        if "case \"$pm\" in" in line:
            skip_pm_translate = False
            new_lines.append(line)
        continue
        
    if "install_qdrant() {" in line:
        new_lines.append(line)
        # We handle container logic in Qdrant and Ollama install functions
        qdrant_container = """
  if [[ "$CONTAINER_MODE" == true ]]; then
    log "Deploying Qdrant via Podman Quadlet..."
    pm_install podman
    mkdir -p "$HOME_DIR/.config/containers/systemd"
    cat > "$HOME_DIR/.config/containers/systemd/qdrant.container" << 'EOF'
[Container]
Image=docker.io/qdrant/qdrant:latest
PublishPort=6333:6333
PublishPort=6334:6334
Volume=qdrant_data:/qdrant/storage

[Service]
Restart=always

[Install]
WantedBy=default.target
EOF
    chown -R "$USER_NAME:$USER_NAME" "$HOME_DIR/.config/containers"
    su - "$USER_NAME" -c "systemctl --user daemon-reload; systemctl --user enable --now qdrant" >>"$LOGFILE" 2>&1
    ok "Qdrant container deployed"
    return 0
  fi
"""
        new_lines.append(qdrant_container)
        continue

    if "install_ollama() {" in line:
        new_lines.append(line)
        ollama_container = """
  if [[ "$CONTAINER_MODE" == true ]]; then
    log "Deploying Ollama via Podman Quadlet..."
    pm_install podman
    mkdir -p "$HOME_DIR/.config/containers/systemd"
    cat > "$HOME_DIR/.config/containers/systemd/ollama.container" << 'EOF'
[Container]
Image=docker.io/ollama/ollama:latest
PublishPort=11434:11434
Volume=ollama_data:/root/.ollama
# If GPU: Add --device nvidia.com/gpu=all or similar depending on setup

[Service]
Restart=always

[Install]
WantedBy=default.target
EOF
    chown -R "$USER_NAME:$USER_NAME" "$HOME_DIR/.config/containers"
    su - "$USER_NAME" -c "systemctl --user daemon-reload; systemctl --user enable --now ollama" >>"$LOGFILE" 2>&1
    ok "Ollama container deployed"
    return 0
  fi
"""
        new_lines.append(ollama_container)
        continue
        
    if "deploy_lsfs() {" in line:
        new_lines.append(line)
        # Add phase 3 tools installation
        phase_3_install = """
  # Phase 3 Tool Installations
  log "Installing Phase 3 tools..."
  mkdir -p /usr/local/bin
  
  if [ -f "/Users/shrey/ash-iso/scripts/ash-workspace.sh" ]; then
    cp /Users/shrey/ash-iso/scripts/ash-workspace.sh /usr/local/bin/ash
    chmod +x /usr/local/bin/ash
  fi
  
  if [ -f "/Users/shrey/ash-iso/ai-services/ash_ask.py" ]; then
    cp /Users/shrey/ash-iso/ai-services/ash_ask.py /usr/local/bin/ash-ask
    chmod +x /usr/local/bin/ash-ask
  fi
  
  if [ -f "/Users/shrey/ash-iso/scripts/agi.sh" ]; then
    cp /Users/shrey/ash-iso/scripts/agi.sh /usr/local/bin/agi
    chmod +x /usr/local/bin/agi
  fi

  if [ -f "/Users/shrey/ash-iso/ai-services/ash_clipboard.py" ]; then
    mkdir -p /opt/ash-ai
    cp /Users/shrey/ash-iso/ai-services/ash_clipboard.py /opt/ash-ai/ash_clipboard.py
    chmod +x /opt/ash-ai/ash_clipboard.py
    
    cat > /etc/systemd/user/ash-clipboard.service << 'EOF'
[Unit]
Description=Ash Semantic Clipboard History
After=default.target

[Service]
Type=simple
ExecStart=/usr/bin/env python3 /opt/ash-ai/ash_clipboard.py
Restart=always

[Install]
WantedBy=default.target
EOF
    # For now, just create it; real enable happens on per-user basis.
  fi
  ok "Phase 3 tools installed"
"""
        new_lines.append(phase_3_install)
        continue

    if "main() {" in line:
        # Generate export-ansible function
        export_ansible = """
export_ansible_role() {
  local out="/tmp/ash-ansible-role"
  mkdir -p "$out"/{tasks,handlers,defaults}
  
  cat > "$out/defaults/main.yml" << 'EOF'
---
qdrant_port: 6333
ollama_port: 11434
model_name: nomic-embed-text
EOF

  cat > "$out/tasks/main.yml" << 'EOF'
---
- name: Install Podman
  package:
    name: podman
    state: present

- name: Create Quadlet directories
  file:
    path: "~/.config/containers/systemd"
    state: directory
    mode: '0755'

- name: Qdrant Quadlet
  copy:
    dest: "~/.config/containers/systemd/qdrant.container"
    content: |
      [Container]
      Image=docker.io/qdrant/qdrant:latest
      PublishPort={{ qdrant_port }}:{{ qdrant_port }}
      Volume=qdrant_data:/qdrant/storage
      [Service]
      Restart=always
      [Install]
      WantedBy=default.target

- name: Ollama Quadlet
  copy:
    dest: "~/.config/containers/systemd/ollama.container"
    content: |
      [Container]
      Image=docker.io/ollama/ollama:latest
      PublishPort={{ ollama_port }}:{{ ollama_port }}
      Volume=ollama_data:/root/.ollama
      [Service]
      Restart=always
      [Install]
      WantedBy=default.target

- name: Daemon reload & start
  command: systemctl --user daemon-reload

- name: Enable Qdrant
  command: systemctl --user enable --now qdrant

- name: Enable Ollama
  command: systemctl --user enable --now ollama

- name: Pull Model
  command: curl -sf -X POST http://localhost:11434/api/pull -d '{"model":"{{ model_name }}"}'
  register: model_pull
EOF

  echo "Ansible role exported to $out"
  echo "Usage: Create a playbook and use this role."
  exit 0
}
"""
        new_lines.append(export_ansible)
        new_lines.append(line)
        continue
        
    if "check_root \"$@\"" in line:
        new_lines.append("  if [[ \"$EXPORT_ANSIBLE\" == true ]]; then export_ansible_role; fi\n")
        new_lines.append(line)
        continue
        
    new_lines.append(line)

with open(INSTALLER_PATH, "w") as f:
    f.writelines(new_lines)

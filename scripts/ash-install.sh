#!/usr/bin/env bash
set -euo pipefail

VERSION="3.0.0"

HOME_DIR="${HOME:-/home/$(whoami)}"
USER_NAME="${USER:-$(whoami)}"
LOG_DIR="/tmp/ash-install"
mkdir -p "$LOG_DIR"
LOGFILE="$LOG_DIR/ash-install-$(date +%Y%m%d-%H%M%S).log"
STATE_FILE="$LOG_DIR/.install-state"
INSTALL_MARK="/etc/ash-installed"

exec 2> >(tee -a "$LOGFILE" >&2)

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; BOLD='\033[1m'; NC='\033[0m'
PASS="${GREEN}✔${NC}"; FAIL="${RED}✘${NC}"; WARN="${YELLOW}⚠${NC}"

declare -a PHASES=()
declare -a PHASE_STATUS=()
ROLLBACK_SNAPSHOT=""

cleanup() { local r=$?; [[ $r -ne 0 ]] && rollback; exit $r; }
trap cleanup EXIT INT TERM

log()  { echo -e "${CYAN}[$(date +%H:%M:%S)]${NC} $*" | tee -a "$LOGFILE"; }
ok()   { echo -e "  ${PASS}  $*" | tee -a "$LOGFILE"; }
fail() { echo -e "  ${FAIL}  $*" | tee -a "$LOGFILE"; }
warn() { echo -e "  ${WARN}  $*" | tee -a "$LOGFILE"; }

phase_start() { PHASES+=("$1"); PHASE_STATUS+=("pending"); log "────────── $1 ──────────"; }
phase_ok()    { local i=$((${#PHASES[@]}-1)); PHASE_STATUS[$i]="pass"; ok "$1"; }
phase_fail()  { local i=$((${#PHASES[@]}-1)); PHASE_STATUS[$i]="fail"; fail "$1"; }

whiptail_installed() { command -v whiptail &>/dev/null; }
dialog_installed()   { command -v dialog &>/dev/null; }
have_tui()           { whiptail_installed || dialog_installed; }

tui_msgbox() {
  local title="$1" msg="$2"
  if whiptail_installed; then whiptail --title "$title" --msgbox "$msg" 20 70 2>/dev/null
  elif dialog_installed; then dialog --title "$title" --msgbox "$msg" 20 70 2>/dev/null
  else echo -e "\n$title\n$msg\n"; fi
}

tui_yesno() {
  local title="$1" msg="$2"
  if whiptail_installed; then whiptail --title "$title" --yesno "$msg" 20 70 2>/dev/null
  elif dialog_installed; then dialog --title "$title" --yesno "$msg" 20 70 2>/dev/null
  else echo -e "\n$title\n$msg\n[Y/n] "; read -r r; [[ "$r" =~ ^[Yy]?$ ]]; fi
}

tui_menu() {
  local title="$1" msg="$2"; shift 2
  local items=("$@")
  if whiptail_installed; then
    local cmd=(whiptail --title "$title" --menu "$msg" 24 70 14)
    for ((i=0; i<${#items[@]}; i+=2)); do cmd+=("${items[$i]}" "${items[$i+1]}"); done
    "${cmd[@]}" 2>/dev/null
  elif dialog_installed; then
    local cmd=(dialog --title "$title" --menu "$msg" 24 70 14)
    for ((i=0; i<${#items[@]}; i+=2)); do cmd+=("${items[$i]}" "${items[$i+1]}"); done
    "${cmd[@]}" 2>/dev/null
  else
    echo "$msg"; for ((i=0; i<${#items[@]}; i+=2)); do echo "  ${items[$i]}) ${items[$i+1]}"; done
    read -r r; echo "$r"
  fi
}

tui_input() {
  local title="$1" msg="$2" default="$3"
  if whiptail_installed; then whiptail --title "$title" --inputbox "$msg" 12 70 "$default" 2>/dev/null
  elif dialog_installed; then dialog --title "$title" --inputbox "$msg" 12 70 "$default" 2>/dev/null
  else echo "$default"; fi
}

tui_password() {
  local title="$1" msg="$2"
  if whiptail_installed; then whiptail --title "$title" --passwordbox "$msg" 10 70 2>/dev/null
  elif dialog_installed; then dialog --title "$title" --passwordbox "$msg" 10 70 2>/dev/null
  else echo "aiuser"; fi
}

tui_gauge() {
  local title="$1" msg="$2" pct="$3"
  if whiptail_installed; then echo "$pct" | whiptail --title "$title" --gauge "$msg" 8 70 0 2>/dev/null || true
  elif dialog_installed; then echo "$pct" | dialog --title "$title" --gauge "$msg" 8 70 0 2>/dev/null || true
  fi
}

redact() { echo "$1" | sed 's/./*/g'; }

run_phase() {
  local name="$1"; shift
  phase_start "$name"
  tui_gauge "Ash Installer" "$name" $(( ${#PHASES[@]} * 100 / 8 ))
  if "$@"; then phase_ok "$name completed"; else phase_fail "$name failed"; return 1; fi
}

check_root() {
  if [[ $EUID -ne 0 ]]; then
    if command -v sudo &>/dev/null; then
      exec sudo HOME="$HOME" bash "$0" "$@"
    else
      echo "Must run as root. Try: su -c 'bash $0'"
      exit 1
    fi
  fi
}

detect_system() {
  if [[ ! -f /etc/arch-release ]]; then
    tui_msgbox "Unsupported System" "Ash Linux is designed for Arch Linux.\nDetected: $(cat /etc/os-release 2>/dev/null | grep '^ID=' | cut -d= -f2 || echo 'unknown')\n\nContinuing may cause issues."
  fi
}

ensure_packages() {
  local pkgs=("$@")
  local missing=()
  for pkg in "${pkgs[@]}"; do
    if ! pacman -Qi "$pkg" &>/dev/null; then
      missing+=("$pkg")
    fi
  done
  if [[ ${#missing[@]} -gt 0 ]]; then
    log "Installing: ${missing[*]}"
    pacman -S --noconfirm --needed "${missing[@]}" >>"$LOGFILE" 2>&1 || return 1
  fi
  return 0
}

ensure_btrfs() {
  local root_dev
  root_dev=$(findmnt -n -o SOURCE / 2>/dev/null || echo "")
  if [[ -z "$root_dev" ]]; then return 1; fi
  if ! mount | grep -q "on / type btrfs"; then return 1; fi
  return 0
}

create_snapshot() {
  if ! ensure_btrfs; then warn "Not on btrfs — skipping snapshots"; return 1; fi
  if ! command -v snapper &>/dev/null; then
    warn "snapper not installed — skipping snapshots"
    return 1
  fi
  local snap_name="ash-install-$(date +%Y%m%d-%H%M%S)"
  log "Creating btrfs snapshot '$snap_name'..."
  if snapper -c root create -d "$snap_name" >>"$LOGFILE" 2>&1; then
    ROLLBACK_SNAPSHOT="$snap_name"
    ok "Snapshot: $snap_name"
    return 0
  else
    warn "Could not create snapper snapshot"
    if which btrfs &>/dev/null; then
      mkdir -p /.snapshots
      btrfs subvolume snapshot / "/.snapshots/$snap_name" >>"$LOGFILE" 2>&1 && {
        ROLLBACK_SNAPSHOT="/.snapshots/$snap_name"
        ok "Raw btrfs snapshot: $ROLLBACK_SNAPSHOT"
        return 0
      }
    fi
    warn "No rollback point created"
    return 1
  fi
}

rollback() {
  if [[ -n "$ROLLBACK_SNAPSHOT" ]]; then
    warn "ROLLBACK: Restoring snapshot $ROLLBACK_SNAPSHOT"
    if snapper -c root list 2>/dev/null | grep -q "$ROLLBACK_SNAPSHOT"; then
      snapper -c root undochange "$ROLLBACK_SNAPSHOT..0" >>"$LOGFILE" 2>&1 || true
    elif [[ -d "$ROLLBACK_SNAPSHOT" ]]; then
      btrfs subvolume list / | grep -q "$(basename "$ROLLBACK_SNAPSHOT")" && {
        btrfs subvolume delete "$ROLLBACK_SNAPSHOT" >>"$LOGFILE" 2>&1 || true
      }
    fi
    ROLLBACK_SNAPSHOT=""
  fi
}

save_state() {
  local phase="$1" status="$2"
  echo "$phase:$status:$(date +%s)" >> "$STATE_FILE"
}

load_state() {
  [[ -f "$STATE_FILE" ]] && grep "$1" "$STATE_FILE" | tail -1 | cut -d: -f2 || echo ""
}

mark_installed() {
  echo "ash-install v$VERSION $(date)" > "$INSTALL_MARK"
  echo "Components:" >> "$INSTALL_MARK"
  for i in "${!PHASES[@]}"; do
    echo "  ${PHASES[$i]}: ${PHASE_STATUS[$i]}" >> "$INSTALL_MARK"
  done
  chmod 644 "$INSTALL_MARK"
}

################################################################################
# COMPONENT INSTALLERS
################################################################################

install_packages() {
  ensure_packages curl wget wofi swaync jq fd openssh kitty python python-pip \
    python-inotify_simple python-requests python-aiohttp lm_sensors || return 1
}

install_qdrant() {
  if curl -sf http://localhost:6333/healthz >/dev/null 2>&1; then
    ok "Qdrant already running"
    return 0
  fi
  if pacman -Qi qdrant &>/dev/null 2>&1; then
    systemctl enable --now qdrant >>"$LOGFILE" 2>&1 || true
    return 0
  fi
  log "Downloading Qdrant binary..."
  local tag ver url tmp
  tag=$(curl -sfL "https://api.github.com/repos/qdrant/qdrant/releases/latest" | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": "\([^"]*\)".*/\1/') || tag="v1.13.6"
  ver="${tag#v}"
  tmp=$(mktemp -d)

  for variant in "musl" "gnu"; do
    url="https://github.com/qdrant/qdrant/releases/download/${tag}/qdrant-x86_64-unknown-linux-${variant}.tar.gz"
    if curl -fSL --progress-bar -o "$tmp/qdrant.tar.gz" "$url" >>"$LOGFILE" 2>&1; then
      break
    fi
  done

  tar -xzf "$tmp/qdrant.tar.gz" -C "$tmp" >>"$LOGFILE" 2>&1
  local bin
  bin=$(find "$tmp" -type f -name "qdrant" | head -1)
  if [[ -z "$bin" ]]; then rm -rf "$tmp"; return 1; fi
  install -m 0755 "$bin" /usr/local/bin/qdrant
  rm -rf "$tmp"

  id -u qdrant &>/dev/null || useradd -r -s /usr/bin/nologin -d /var/lib/qdrant qdrant
  mkdir -p /var/lib/qdrant
  chown -R qdrant:qdrant /var/lib/qdrant

  cat > /etc/systemd/system/qdrant.service << 'UNIT'
[Unit]
Description=Qdrant Vector Search Engine
After=network-online.target
Wants=network-online.target
[Service]
Type=simple
User=qdrant
Group=qdrant
ExecStart=/usr/local/bin/qdrant --storage-path /var/lib/qdrant
Restart=always
RestartSec=2
LimitNOFILE=65536
Environment=QDRANT__SERVICE__HTTP_PORT=6333
[Install]
WantedBy=multi-user.target
UNIT

  systemctl daemon-reload >>"$LOGFILE" 2>&1
  systemctl enable --now qdrant >>"$LOGFILE" 2>&1

  for i in $(seq 1 30); do
    if curl -sf --max-time 2 "http://localhost:6333/healthz" >/dev/null 2>&1; then
      ok "Qdrant healthy on :6333"
      return 0
    fi
    sleep 1
  done
  return 1
}

install_ollama() {
  if curl -sf http://localhost:11434/api/version >/dev/null 2>&1; then
    ok "Ollama already running"
    return 0
  fi
  if command -v ollama &>/dev/null; then
    systemctl enable --now ollama >>"$LOGFILE" 2>&1 || true
  else
    log "Installing Ollama..."
    if pacman -Qi ollama &>/dev/null 2>&1; then
      systemctl enable --now ollama >>"$LOGFILE" 2>&1 || true
    else
      curl -sfL https://ollama.com/install.sh >>"$LOGFILE" 2>&1 | sh 2>>"$LOGFILE" || {
        warn "Ollama install script failed — trying AUR..."
        if command -v paru &>/dev/null; then
          su - "$USER_NAME" -c "paru -S --noconfirm ollama" >>"$LOGFILE" 2>&1 || true
        fi
      }
      sleep 3
      systemctl enable --now ollama >>"$LOGFILE" 2>&1 || true
    fi
  fi
  for i in $(seq 1 30); do
    if curl -sf --max-time 2 "http://localhost:11434/api/version" >/dev/null 2>&1; then
      ok "Ollama running on :11434"
      return 0
    fi
    sleep 1
  done
  return 1
}

pull_model() {
  local model="${1:-nomic-embed-text}"
  if curl -sf http://localhost:11434/api/tags 2>/dev/null | grep -q "$model"; then
    ok "Model $model already present"
    return 0
  fi
  log "Pulling $model (this can take a while) ..."
  su - "$USER_NAME" -c "ollama pull $model" >>"$LOGFILE" 2>&1 || {
    curl -sf -X POST http://localhost:11434/api/pull \
      -d "{\"model\":\"$model\"}" >>"$LOGFILE" 2>&1 || return 1
  }
  curl -sf -X POST http://localhost:11434/api/generate \
    -d "{\"model\":\"$model\",\"keep_alive\":-1,\"prompt\":\"warmup\"}" -o /dev/null >>"$LOGFILE" 2>&1 || true
  ok "Model $model ready"
}

deploy_lsfs() {
  mkdir -p "$HOME_DIR/.config/scripts" "$HOME_DIR/.local/bin" "$HOME_DIR/.config/systemd/user"
  chown -R "$USER_NAME:$USER_NAME" "$HOME_DIR/.config" "$HOME_DIR/.local" 2>/dev/null || true

  local daemon_py="$HOME_DIR/.config/scripts/lsfs_daemon.py"
  local launcher_sh="$HOME_DIR/.config/scripts/lsfs_launcher_hook.sh"
  local query_py="$HOME_DIR/.config/scripts/lsfs_query.py"
  local svc="$HOME_DIR/.config/systemd/user/lsfs-daemon.service"

  cat > "$query_py" << 'PYQUERY'
#!/usr/bin/env python3
import sys, json, os, urllib.request, http.client, socket

OLLAMA_URL = "http://localhost:11434/api/embeddings"
MODEL = "nomic-embed-text"
QDRANT_TCP = "http://localhost:6333"
COSINE_FLOOR = 0.5
SEARCH_LIMIT = 20

def qdrant_req(method, path, data=None, timeout=2):
    api_path = f"/collections/apps/{path.lstrip('/')}"
    body = json.dumps(data).encode() if data else None
    headers = {"Content-Type": "application/json"}
    url = f"{QDRANT_TCP}{api_path}"
    req = urllib.request.Request(url, data=body, method=method, headers=headers)
    try:
        with urllib.request.urlopen(req, timeout=timeout) as r:
            return json.loads(r.read())
    except Exception:
        return {}

def search(query, limit=SEARCH_LIMIT):
    req_data = json.dumps({"model": MODEL, "prompt": query[:2048], "keep_alive": -1}).encode()
    req = urllib.request.Request(OLLAMA_URL, data=req_data, headers={"Content-Type": "application/json"})
    with urllib.request.urlopen(req, timeout=3.0) as resp:
        emb_data = json.loads(resp.read())
    embedding = emb_data.get("embedding")
    if not embedding:
        return []
    search_payload = {"vector": embedding, "limit": limit, "with_payload": True}
    hits = qdrant_req("POST", "points/search", search_payload, timeout=2.0)
    results = []
    seen = set()
    for hit in hits.get("result", []):
        score = hit.get("score", 0)
        payload = hit.get("payload", {})
        path = payload.get("path", "")
        if score < COSINE_FLOOR or path in seen:
            continue
        seen.add(path)
        results.append({"path": path, "name": payload.get("name", ""), "score": score})
    results.sort(key=lambda r: r['score'], reverse=True)
    return results

if __name__ == "__main__":
    query = " ".join(sys.argv[1:]) if len(sys.argv) > 1 else ""
    if not query:
        print("Usage: lsfs-query <query>")
        sys.exit(0)
    results = search(query)
    if not results:
        print("No matches found.")
    for r in results:
        print(f"{r['path']} | {r['name']} ({r['score']:.3f})")
PYQUERY

  cat > "$daemon_py" << 'PYDAEMON'
#!/usr/bin/env python3
import os, sys, json, time, hashlib, logging, requests
from pathlib import Path

HOME = os.environ.get("HOME", "/home/pal")
OLLAMA = "http://localhost:11434/api/embeddings"
QDRANT = "http://localhost:6333/collections"
MODEL = "nomic-embed-text"
COLLECTION = "apps"
WATCH = [HOME]
IGNORE = {".git", "node_modules", "__pycache__", ".cache", ".venv", "venv", ".mozilla", ".steam"}

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
log = logging.getLogger("lsfs")

def ensure_collection():
    r = requests.get(f"{QDRANT}/{COLLECTION}", timeout=5)
    if r.status_code == 200:
        return
    payload = {"vectors": {"size": 768, "distance": "Cosine"}}
    for i in range(5):
        r = requests.put(f"{QDRANT}/{COLLECTION}", json=payload, timeout=10)
        if r.status_code in (200, 201):
            log.info("Collection created")
            return
        time.sleep(2)
    log.error("Could not create Qdrant collection")

def embed(text):
    for i in range(3):
        r = requests.post(OLLAMA, json={"model": MODEL, "prompt": text[:512], "keep_alive": -1}, timeout=30)
        if r.status_code == 200:
            data = r.json()
            if "embedding" in data:
                return data["embedding"]
        time.sleep(2)
    return None

def index_file(path):
    for p in IGNORE:
        if p in path:
            return False
    fp = Path(path)
    if not fp.is_file() or fp.is_symlink():
        return False
    try:
        stat = fp.stat()
        if stat.st_size == 0 or stat.st_size > 1048576:
            return False
        text = fp.read_text(errors="replace")[:1024]
        if not text.strip():
            return False
        vec = embed(text)
        if not vec:
            return False
        pid = hashlib.md5(path.encode()).hexdigest()
        payload = {"path": path, "name": fp.name, "ext": fp.suffix.lower(), "mtime": int(stat.st_mtime), "size": stat.st_size}
        r = requests.put(f"{QDRANT}/{COLLECTION}/points", json={
            "points": [{"id": pid, "vector": vec, "payload": payload}]
        }, timeout=10)
        return r.status_code in (200, 201)
    except Exception:
        return False

def full_scan():
    log.info("Full scan started")
    indexed = 0
    for root, dirs, files in os.walk(HOME):
        dirs[:] = [d for d in dirs if d not in IGNORE and not d.startswith(".")]
        for f in files:
            try:
                if index_file(os.path.join(root, f)):
                    indexed += 1
            except Exception:
                pass
    log.info(f"Full scan done: {indexed} files indexed")

def watch_loop():
    log.info("Polling every 60s for changes")
    known = set()
    while True:
        current = set()
        for root, dirs, files in os.walk(HOME):
            dirs[:] = [d for d in dirs if d not in IGNORE and not d.startswith(".")]
            for f in files:
                fp = os.path.join(root, f)
                current.add(fp)
                if fp not in known:
                    index_file(fp)
        known = current
        time.sleep(60)

if __name__ == "__main__":
    log.info("LSFS Daemon starting")
    for i in range(60):
        if requests.get(f"{QDRANT}/{COLLECTION}", timeout=2).status_code == 200:
            break
        time.sleep(1)
    else:
        log.error("Qdrant not reachable")
        sys.exit(1)
    ensure_collection()
    full_scan()
    watch_loop()
PYDAEMON

  cat > "$launcher_sh" << 'LSFSHOOK'
#!/usr/bin/env bash
set -euo pipefail
QUERY=$(wofi --dmenu --prompt "Agentic Search" --cache-file /dev/null < /dev/null 2>/dev/null) || true
[ -z "${QUERY:-}" ] && exit 0
notify-send -t 3000 -r 999 "Agentic OS" "Searching: $QUERY"
RESULTS_FILE=$(mktemp /tmp/lsfs_results.XXXXXX)
trap 'rm -f "$RESULTS_FILE"' EXIT
HAS_RESULTS=0

if command -v curl &>/dev/null; then
  OLLAMA_RESP=$(curl -s --max-time 5 -X POST http://localhost:11434/api/embeddings \
    -d "{\"model\":\"nomic-embed-text\",\"prompt\":\"${QUERY//\"/\\\"}\",\"keep_alive\":-1}" 2>/dev/null) || true
  if [ -n "$OLLAMA_RESP" ]; then
    EMBEDDING=$(echo "$OLLAMA_RESP" | tr -d '\n' | sed -n 's/.*"embedding":\(\[[^]]*\]\).*/\1/p') || true
    if [ -n "${EMBEDDING:-}" ]; then
      VECTOR=$(echo "$EMBEDDING" | tr -d ' \t\n')
      QDRANT_RESP=$(curl -s --max-time 5 -X POST http://localhost:6333/collections/apps/points/search \
        -H "Content-Type: application/json" \
        -d "{\"vector\":$VECTOR,\"limit\":10,\"with_payload\":true}" 2>/dev/null) || true
      if [ -n "${QDRANT_RESP:-}" ]; then
        echo "$QDRANT_RESP" | python3 -c "
import sys, json
data = json.load(sys.stdin)
for r in data.get('result', []):
    p = r.get('payload', {})
    print(f\"{p.get('path', '')} | {p.get('name', '')} ({r.get('score', 0):.3f})\")
" > "$RESULTS_FILE" 2>/dev/null
        [ -s "$RESULTS_FILE" ] && HAS_RESULTS=1
      fi
    fi
  fi
fi

if [ "$HAS_RESULTS" -eq 0 ]; then
  PAT=$(echo "$QUERY" | grep -oiE '[0-9]+\s*(h|hr|hour|hours|d|day|days)' | head -1 | tr -d ' ') || true
  if [ -n "$PAT" ]; then
    if echo "$PAT" | grep -qiE '[0-9]+[hd]'; then TIME_ARG="$PAT"
    else
      NUM=$(echo "$PAT" | grep -oE '[0-9]+')
      echo "$PAT" | grep -qiE 'h|hr|hour|hours' && TIME_ARG="${NUM}h" || TIME_ARG="${NUM}d"
    fi
    if command -v fd &>/dev/null; then fd --changed-within "$TIME_ARG" --type f "$HOME" 2>/dev/null | head -20 > "$RESULTS_FILE" || true
    elif command -v find &>/dev/null; then
      echo "$TIME_ARG" | grep -q 'h$' && MINS=$(( ${TIME_ARG%h} * 60 )) || MINS=$(( ${TIME_ARG%d} * 1440 ))
      find "$HOME" -mmin "-${MINS}" -type f 2>/dev/null | head -20 > "$RESULTS_FILE" || true
    fi
    [ -s "$RESULTS_FILE" ] && HAS_RESULTS=1
  fi
fi

if [ "$HAS_RESULTS" -eq 0 ] && command -v fd &>/dev/null; then
  fd --type f --max-depth 5 "$HOME" 2>/dev/null | head -15 > "$RESULTS_FILE" || true
  [ -s "$RESULTS_FILE" ] && HAS_RESULTS=1
fi

if [ "$HAS_RESULTS" -eq 0 ]; then notify-send -u critical -t 5000 "Agentic OS" "No files found"; exit 0; fi

notify-send -t 2000 -r 999 "Agentic OS" "$(wc -l < "$RESULTS_FILE") results"
SELECTED=$(wofi --dmenu --prompt "Results ($(wc -l < "$RESULTS_FILE"))" --cache-file /dev/null < "$RESULTS_FILE" 2>/dev/null) || true
[ -z "${SELECTED:-}" ] && exit 0
TARGET_PATH=$(echo "$SELECTED" | sed 's/ | .*//; s/\t.*//')
[ -z "$TARGET_PATH" ] && exit 0
[ ! -e "$TARGET_PATH" ] && notify-send -u critical "Not found: $TARGET_PATH" && exit 1
if [ -d "$TARGET_PATH" ]; then hyprctl dispatch exec "kitty -e yazi '$TARGET_PATH'" 2>/dev/null || true
elif echo "$TARGET_PATH" | grep -q '\.desktop$'; then hyprctl dispatch exec "gtk-launch '$(basename "$TARGET_PATH")'" 2>/dev/null || true
else hyprctl dispatch exec "kitty --class floating_editor -e nvim '$TARGET_PATH'" 2>/dev/null || true; fi
LSFSHOOK

  chmod +x "$daemon_py" "$launcher_sh" "$query_py"
  chown -R "$USER_NAME:$USER_NAME" "$HOME_DIR/.config/scripts"

  for name in lsfs_daemon lsfs_query; do
    ln -sf "$HOME_DIR/.config/scripts/${name}.py" "$HOME_DIR/.local/bin/${name/_/-}" 2>/dev/null || true
  done

  for rc in bashrc zshrc; do
    local rcpath="$HOME_DIR/.$rc"
    if [[ -f "$rcpath" ]] && ! grep -q '\.local/bin' "$rcpath" 2>/dev/null; then
      echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$rcpath"
    fi
  done

  cat > "$HOME_DIR/.lsfsignore" << 'IGNORE'
.*
node_modules/ __pycache__/ *.pyc .env/ venv/ .venv/ target/ build/ dist/
*.egg-info/ site-packages/ .git/ .svn/
*.csv *.log *.sql *.db *.sqlite *.pkl
*.mp3 *.mp4 *.png *.jpg *.ico *.svg
*.zip *.tar *.gz *.xz *.bz2 *.rar *.7z
*.o *.so *.dylib *.dll *.exe *.bin
.idea/ .vscode/ *.swp *~ .DS_Store
lost+found/ .Trash/
IGNORE
  chown "$USER_NAME:$USER_NAME" "$HOME_DIR/.lsfsignore"

  cat > "$svc" << 'SYSD'
[Unit]
Description=LSFS Semantic File Indexer
After=network.target
[Service]
Type=simple
ExecStart=%h/.config/scripts/lsfs_daemon.py
Restart=always
RestartSec=5
Nice=19
IOSchedulingClass=idle
CPUQuota=30%
Environment=PYTHONUNBUFFERED=1
[Install]
WantedBy=default.target
SYSD
  chown "$USER_NAME:$USER_NAME" "$svc"

  loginctl enable-linger "$USER_NAME" >>"$LOGFILE" 2>&1 || true
  su - "$USER_NAME" -c "XDG_RUNTIME_DIR=/run/user/$(id -u $USER_NAME) systemctl --user daemon-reload" >>"$LOGFILE" 2>&1 || true
  su - "$USER_NAME" -c "XDG_RUNTIME_DIR=/run/user/$(id -u $USER_NAME) systemctl --user enable --now lsfs-daemon.service" >>"$LOGFILE" 2>&1 || {
    warn "Daemon service failed — direct launch"
    su - "$USER_NAME" -c "nohup python3 $HOME_DIR/.config/scripts/lsfs_daemon.py &>/tmp/lsfs-daemon.log &" >>"$LOGFILE" 2>&1 || true
  }

  return 0
}

patch_hyprland() {
  local hypr_conf="$HOME_DIR/.config/hypr/hyprland.conf"
  if [[ ! -f "$hypr_conf" ]]; then
    mkdir -p "$HOME_DIR/.config/hypr"
    echo "# Ash Linux Hyprland" > "$hypr_conf"
  fi
  if grep -q "lsfs_launcher_hook" "$hypr_conf" 2>/dev/null; then
    ok "Super+Space already configured"
  else
    cat >> "$hypr_conf" << 'EOF'
bind = SUPER, Space, exec, $HOME/.config/scripts/lsfs_launcher_hook.sh
EOF
    ok "Super+Space bound to launcher"
  fi
  chown -R "$USER_NAME:$USER_NAME" "$HOME_DIR/.config/hypr" 2>/dev/null || true
  su - "$USER_NAME" -c "hyprctl reload" 2>/dev/null || true
}

enable_auto_update() {
  cat > /etc/systemd/system/ash-auto-update.service << 'UNIT'
[Unit]
Description=Ash Linux Auto-Update
[Service]
Type=oneshot
ExecStart=/bin/bash -c "curl -sfL https://ash.sh/install | bash"
Nice=19
IOSchedulingClass=idle
UNIT
  cat > /etc/systemd/system/ash-auto-update.timer << 'TIMER'
[Unit]
Description=Weekly Ash Linux update
[Timer]
OnCalendar=weekly
Persistent=true
RandomizedDelaySec=2h
[Install]
WantedBy=timers.target
TIMER
  systemctl daemon-reload >>"$LOGFILE" 2>&1
  systemctl enable ash-auto-update.timer >>"$LOGFILE" 2>&1 || true
}

print_dashboard() {
  echo ""
  echo -e "${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
  echo -e "${GREEN}║  Ash Linux v$VERSION — Installation Complete                    ║${NC}"
  echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"
  echo ""

  printf "  ${BOLD}%-28s %s${NC}\n" "Component" "Status"
  printf "  %-28s %s\n" "────────────────────────────" "──────────"
  local all_pass=true
  for i in "${!PHASES[@]}"; do
    local name="${PHASES[$i]}"
    local result="${PHASE_STATUS[$i]}"
    if [[ "$result" == "pass" ]]; then
      printf "  %-28s ${GREEN}%-10s${NC}\n" "$name" "✔ PASS"
    else
      printf "  %-28s ${RED}%-10s${NC}\n" "$name" "✘ FAIL"
      all_pass=false
    fi
  done

  echo ""
  local q=0 o=0 d=0 l=0
  curl -sf http://localhost:6333/healthz >/dev/null 2>&1 && q=1
  curl -sf http://localhost:11434/api/version >/dev/null 2>&1 && o=1
  pgrep -f lsfs_daemon >/dev/null 2>&1 && d=1
  [[ -x "$HOME_DIR/.config/scripts/lsfs_launcher_hook.sh" ]] && l=1

  echo -e "  ${CYAN}Live Status:${NC}"
  echo -e "    Qdrant:    $([ $q -eq 1 ] && echo "${GREEN}Running${NC}" || echo "${RED}Down${NC}")"
  echo -e "    Ollama:    $([ $o -eq 1 ] && echo "${GREEN}Running${NC}" || echo "${RED}Down${NC}")"
  echo -e "    LSFS:      $([ $d -eq 1 ] && echo "${GREEN}Running${NC}" || echo "${RED}Down${NC}")"
  echo -e "    Launcher:  $([ $l -eq 1 ] && echo "${GREEN}Ready${NC}" || echo "${RED}Missing${NC}")"

  echo ""
  echo -e "  ${BOLD}Quick Start:${NC}"
  echo -e "    ${CYAN}Super+Space${NC}  →  Search files by meaning"
  echo -e "    ${CYAN}lsfs-query${NC}    →  CLI search (e.g. lsfs-query 'my config files')"
  echo ""
  echo -e "  ${BOLD}Monitoring:${NC}"
  echo -e "    ${CYAN}journalctl --user -u lsfs-daemon -f${NC}  →  Live indexer logs"
  echo -e "    ${CYAN}ash-doctor${NC}                           →  Full system health"
  echo ""
  echo -e "  ${BOLD}Commands:${NC}"
  echo -e "    ${CYAN}sudo systemctl status qdrant${NC}         →  Vector DB status"
  echo -e "    ${CYAN}sudo systemctl status ollama${NC}         →  Embedding model status"
  echo -e "    ${CYAN}systemctl --user status lsfs-daemon${NC}  →  Indexer status"
  echo ""

  if [[ -n "$ROLLBACK_SNAPSHOT" ]]; then
    echo -e "  ${WARN} Rollback snapshot: $ROLLBACK_SNAPSHOT"
    echo ""
  fi

  if $all_pass; then
    echo -e "  ${GREEN}${BOLD}All systems operational. Press Super+Space to search!${NC}"
  else
    echo -e "  ${YELLOW}${BOLD}Some components need attention. Run ash-doctor to diagnose.${NC}"
  fi

  echo -e "  ${BOLD}Log:${NC} $LOGFILE"
  echo ""
}

################################################################################
# MAIN
################################################################################

main() {
  clear
  echo -e "${CYAN}"
  cat << 'EOF'
    █████╗ ███████╗██╗  ██╗
   ██╔══██╗██╔════╝██║  ██║
   ███████║███████╗███████║
   ██╔══██║╚════██║██╔══██║
   ██║  ██║███████║██║  ██║
   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
EOF
  echo -e "${NC}"
  echo -e "  Ash Linux Installer v$VERSION"
  echo -e "  Single-run setup — Qdrant + Ollama + LSFS + Launcher"
  echo ""

  check_root "$@"
  detect_system

  tui_msgbox "Ash Linux v$VERSION" "This installer will configure:\n\n  1. System packages\n  2. Qdrant vector database\n  3. Ollama AI engine\n  4. Embedding model (nomic-embed-text)\n  5. LSFS semantic indexer\n  6. Super+Space launcher\n  7. Desktop integration\n  8. Auto-updates\n\nA btrfs snapshot will be created for rollback."
  tui_yesno "Begin Installation?" "This will install and configure the full Ash Linux AI stack.\n\nProceed?" || {
    echo "Aborted."
    exit 0
  }

  create_snapshot

  run_phase "System Packages" install_packages || true
  run_phase "Qdrant Vector DB" install_qdrant || true
  run_phase "Ollama AI Engine" install_ollama || true
  run_phase "Embedding Model" pull_model nomic-embed-text || true
  run_phase "LSFS Indexer" deploy_lsfs || true
  run_phase "Desktop Launcher" patch_hyprland || true
  run_phase "Auto Updates" enable_auto_update || true
  run_phase "Final Verification" verify_all || true

  mark_installed
  print_dashboard

  tui_msgbox "Ash Linux Installed" "Press Super+Space to search your files by meaning!\n\nCLI: lsfs-query 'your query'\nLogs: journalctl --user -u lsfs-daemon -f"
}

verify_all() {
  local ok=true
  curl -sf http://localhost:6333/healthz >/dev/null 2>&1 && ok "Qdrant: healthy" || { fail "Qdrant: down"; ok=false; }
  curl -sf http://localhost:11434/api/version >/dev/null 2>&1 && ok "Ollama: running" || { fail "Ollama: down"; ok=false; }
  curl -sf http://localhost:11434/api/tags 2>/dev/null | grep -q nomic-embed-text && ok "Model: loaded" || { warn "Model: not found"; }
  su - "$USER_NAME" -c "XDG_RUNTIME_DIR=/run/user/$(id -u $USER_NAME) systemctl --user is-active lsfs-daemon.service" 2>/dev/null | grep -qE "active|activating" && ok "LSFS daemon: active" || { warn "LSFS daemon: inactive"; }
  if [[ -x "$HOME_DIR/.config/scripts/lsfs_launcher_hook.sh" ]]; then
    ok "Launcher hook: ready"
  else
    fail "Launcher hook: missing"; ok=false
  fi
  if [[ -f "/etc/systemd/system/ash-auto-update.timer" ]]; then
    ok "Auto-update: configured"
  fi
  if which wofi &>/dev/null; then
    ok "Wofi: available"
  fi
  $ok
}

main "$@"

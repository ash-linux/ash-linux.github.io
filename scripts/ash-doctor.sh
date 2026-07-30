#!/usr/bin/env bash
# ash-doctor - Comprehensive Diagnostic Tool for Ash Linux
# Categories: Services, Storage, Network, Security, Performance, Config

set -eo pipefail

JSON_MODE=false
HTML_MODE=false
for arg in "$@"; do
    if [ "$arg" == "--json" ]; then JSON_MODE=true; fi
    if [ "$arg" == "--html" ]; then HTML_MODE=true; fi
    if [ "$arg" == "--audit-report" ]; then
        echo "=== Audit Log Report ==="
        if [ -f "/var/log/ash/audit.log" ]; then
            wc -l /var/log/ash/audit.log
            echo "Recent queries:"
            tail -n 5 /var/log/ash/audit.log
        else
            echo "No audit log found."
        fi
        exit 0
    fi
done

RESULTS=()
FAILURES=0
WARNINGS=0

check() {
    local category="$1"
    local name="$2"
    local cmd="$3"
    local status="PASS"
    local msg=""
    
    if eval "$cmd" > /tmp/ash-doctor-out 2>&1; then
        msg=$(cat /tmp/ash-doctor-out | head -n 1)
        if [ -z "$msg" ]; then msg="OK"; fi
    else
        status="FAIL"
        msg=$(cat /tmp/ash-doctor-out | head -n 1)
        if [ -z "$msg" ]; then msg="Failed"; fi
        FAILURES=$((FAILURES + 1))
    fi
    RESULTS+=("$category|$name|$status|$msg")
}

check_warn() {
    local category="$1"
    local name="$2"
    local cmd="$3"
    local status="PASS"
    local msg=""
    
    if eval "$cmd" > /tmp/ash-doctor-out 2>&1; then
        msg=$(cat /tmp/ash-doctor-out | head -n 1)
        if [ -z "$msg" ]; then msg="OK"; fi
    else
        status="WARN"
        msg=$(cat /tmp/ash-doctor-out | head -n 1)
        if [ -z "$msg" ]; then msg="Warning"; fi
        WARNINGS=$((WARNINGS + 1))
    fi
    RESULTS+=("$category|$name|$status|$msg")
}

# --- SERVICES ---
check "Services" "Qdrant API" "curl -sf http://localhost:6333/healthz"
check "Services" "Ollama API" "curl -sf http://localhost:11434/api/version"
check_warn "Services" "LSFS Daemon" "systemctl --user is-active lsfs-daemon.service 2>/dev/null || echo 'inactive'"

# --- STORAGE ---
check "Storage" "Root Filesystem" "df -h / | tail -n 1 | awk '{print \$5}' | grep -v '100%'"
check "Storage" "Qdrant Directory" "[ -d /var/lib/qdrant ]"
check_warn "Storage" "Btrfs Check" "btrfs filesystem show / >/dev/null 2>&1 || true"

# --- NETWORK ---
check "Network" "Qdrant Port (6333)" "ss -tuln | grep -q 6333"
check "Network" "Ollama Port (11434)" "ss -tuln | grep -q 11434"
check_warn "Network" "Firewall active" "iptables -L >/dev/null 2>&1 || true"

# --- SECURITY ---
check "Security" "AppArmor Enabled" "aa-status >/dev/null 2>&1 || true"
check "Security" "Secrets Permissions" "stat -c '%a' /etc/ash/system.key 2>/dev/null | grep -q '600' || true"
check "Security" "Audit Log Integrity" "lsattr /var/log/ash/audit.log 2>/dev/null | grep -q '\-a\-' || true"

# --- PERFORMANCE ---
check "Performance" "Free RAM" "free -m | awk '/Mem:/ {if (\$4 > 500) exit 0; else exit 1}'"
check_warn "Performance" "Swap" "free -m | awk '/Swap:/ {if (\$2 > 0) exit 0; else exit 1}'"

# Output
if [ "$JSON_MODE" = true ]; then
    echo "{"
    echo "  \"failures\": $FAILURES,"
    echo "  \"warnings\": $WARNINGS,"
    echo "  \"checks\": ["
    FIRST=true
    for r in "${RESULTS[@]}"; do
        IFS='|' read -r cat name stat msg <<< "$r"
        if [ "$FIRST" = true ]; then FIRST=false; else echo ","; fi
        echo "    {\"category\": \"$cat\", \"name\": \"$name\", \"status\": \"$stat\", \"message\": \"$msg\"}" | tr -d '\n'
    done
    echo ""
    echo "  ]"
    echo "}"
elif [ "$HTML_MODE" = true ]; then
    echo "<!DOCTYPE html><html><body><h1>Ash Doctor Report</h1>"
    echo "<p>Failures: $FAILURES, Warnings: $WARNINGS</p><table border='1'>"
    echo "<tr><th>Category</th><th>Check</th><th>Status</th><th>Message</th></tr>"
    for r in "${RESULTS[@]}"; do
        IFS='|' read -r cat name stat msg <<< "$r"
        echo "<tr><td>$cat</td><td>$name</td><td>$stat</td><td>$msg</td></tr>"
    done
    echo "</table></body></html>"
else
    echo "=== Ash Doctor Report ==="
    for r in "${RESULTS[@]}"; do
        IFS='|' read -r cat name stat msg <<< "$r"
        color="\e[32m"
        if [ "$stat" == "FAIL" ]; then color="\e[31m"; fi
        if [ "$stat" == "WARN" ]; then color="\e[33m"; fi
        echo -e "[$cat] $name: ${color}$stat\e[0m ($msg)"
    done
    echo "Failures: $FAILURES | Warnings: $WARNINGS"
fi

if [ "$FAILURES" -gt 0 ]; then
    exit 2
elif [ "$WARNINGS" -gt 0 ]; then
    exit 1
else
    exit 0
fi

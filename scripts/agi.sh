#!/usr/bin/env bash

set -e

# agi.sh
# Implements agi vibecoding features: install, scaffold, sync-ignore

OLLAMA_API="http://localhost:11434/api/generate"
MODEL="llama3.2"

agi_install() {
    local desc="$*"
    if [ -z "$desc" ]; then
        echo "Usage: agi install <natural-language description>"
        exit 1
    fi
    
    echo "Asking Ollama to suggest a package for: $desc"
    local prompt="Suggest the exact package name for the user's OS to install: $desc. Respond with ONLY the package name, nothing else."
    
    local pkg
    pkg=$(curl -s -X POST $OLLAMA_API -d "{\"model\":\"$MODEL\",\"prompt\":\"$prompt\",\"stream\":false}" | jq -r '.response' | tr -d ' ' | tr -d '\n')
    
    if [ -z "$pkg" ]; then
        echo "Failed to get a suggestion from Ollama."
        exit 1
    fi
    
    echo "Suggested package: $pkg"
    read -p "Do you want to install $pkg? [Y/n] " confirm
    if [[ "$confirm" == "" || "$confirm" == "y" || "$confirm" == "Y" ]]; then
        # Detect PM
        if command -v apt >/dev/null; then sudo apt install -y "$pkg"
        elif command -v dnf >/dev/null; then sudo dnf install -y "$pkg"
        elif command -v pacman >/dev/null; then sudo pacman -S --noconfirm "$pkg"
        elif command -v zypper >/dev/null; then sudo zypper install -y "$pkg"
        elif command -v apk >/dev/null; then sudo apk add "$pkg"
        else
            echo "No known package manager found."
            exit 1
        fi
    else
        echo "Aborted."
    fi
}

agi_scaffold() {
    local desc="$*"
    if [ -z "$desc" ]; then
        echo "Usage: agi scaffold <description>"
        exit 1
    fi
    
    echo "Asking Ollama to scaffold project for: $desc"
    local prompt="Provide a bash script to scaffold a project structure for: $desc. Only output the raw bash script."
    
    local script
    script=$(curl -s -X POST $OLLAMA_API -d "{\"model\":\"$MODEL\",\"prompt\":\"$prompt\",\"stream\":false}" | jq -r '.response')
    
    # Strip markdown if present
    script=$(echo "$script" | sed 's/^```bash//; s/^```sh//; s/^```//; s/```$//')
    
    echo "Proposed script:"
    echo "$script"
    echo "-----------------"
    read -p "Execute this script? [y/N] " confirm
    if [[ "$confirm" == "y" || "$confirm" == "Y" ]]; then
        eval "$script"
        echo "Project scaffolded."
    else
        echo "Aborted."
    fi
}

agi_sync_ignore() {
    if [ -f ".gitignore" ]; then
        echo "Syncing .gitignore to .lsfsignore"
        cp .gitignore .lsfsignore
        echo ".lsfsignore updated."
    else
        echo "No .gitignore found in current directory."
    fi
}

case "$1" in
    install)
        shift
        agi_install "$@"
        ;;
    scaffold)
        shift
        agi_scaffold "$@"
        ;;
    sync-ignore)
        agi_sync_ignore
        ;;
    *)
        echo "Usage: agi <install|scaffold|sync-ignore> ..."
        ;;
esac

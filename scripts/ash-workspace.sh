#!/usr/bin/env bash

set -e

# ash-workspace.sh
# Implements Phase 3 vibecoding & developer experience workspace features.

show_help() {
    echo "Usage: ash workspace <command> [<args>]"
    echo "       ash <up|down|status>"
    echo ""
    echo "Commands:"
    echo "  workspace create <name>   Create a new workspace"
    echo "  workspace reset <name>    Reset a workspace to a clean state"
    echo "  workspace list            List all workspaces"
    echo "  workspace delete <name>   Delete a workspace"
    echo "  workspace switch <name>   Switch to a workspace"
    echo "  up                        Start AI services (Qdrant, Ollama, lsfs-daemon)"
    echo "  down                      Stop AI services"
    echo "  status                    Show status of AI services"
}

WORKSPACE_DIR="$HOME/.ash/workspaces"
WORKSPACE_CONF="$HOME/.config/ash/workspace.conf"

create_workspace() {
    local name="$1"
    if [ -z "$name" ]; then echo "Error: Missing workspace name"; exit 1; fi
    mkdir -p "$WORKSPACE_DIR"
    
    if command -v btrfs >/dev/null 2>&1 && btrfs filesystem df "$WORKSPACE_DIR" >/dev/null 2>&1; then
        echo "Creating btrfs subvolume for $name..."
        btrfs subvolume create "$WORKSPACE_DIR/$name"
    else
        echo "Creating directory for $name..."
        mkdir -p "$WORKSPACE_DIR/$name"
    fi
    
    # Dev env symlink
    ln -sfn "$WORKSPACE_DIR/$name" "$HOME/current-workspace"
    echo "Workspace $name created and linked to ~/current-workspace"
}

reset_workspace() {
    local name="$1"
    if [ -z "$name" ]; then echo "Error: Missing workspace name"; exit 1; fi
    local target="$WORKSPACE_DIR/$name"
    
    if [ ! -e "$target" ]; then
        echo "Error: Workspace $name does not exist."
        exit 1
    fi
    
    # In a full implementation we'd use btrfs snapshots, but let's fall back to rm -rf + recreate
    echo "Resetting workspace $name..."
    if command -v btrfs >/dev/null 2>&1 && btrfs subvolume show "$target" >/dev/null 2>&1; then
        btrfs subvolume delete "$target"
        btrfs subvolume create "$target"
    else
        rm -rf "${target:?}"/*
        rm -rf "${target:?}"/.* 2>/dev/null || true
    fi
    echo "Workspace $name reset."
}

list_workspaces() {
    echo "Workspaces in $WORKSPACE_DIR:"
    if [ -d "$WORKSPACE_DIR" ]; then
        ls -lh "$WORKSPACE_DIR" | awk '{print $9, $5, $6, $7, $8}' | grep -v "^$"
    else
        echo "No workspaces found."
    fi
}

delete_workspace() {
    local name="$1"
    if [ -z "$name" ]; then echo "Error: Missing workspace name"; exit 1; fi
    local target="$WORKSPACE_DIR/$name"
    
    if command -v btrfs >/dev/null 2>&1 && btrfs subvolume show "$target" >/dev/null 2>&1; then
        btrfs subvolume delete "$target"
    else
        rm -rf "$target"
    fi
    echo "Workspace $name deleted."
}

switch_workspace() {
    local name="$1"
    if [ -z "$name" ]; then echo "Error: Missing workspace name"; exit 1; fi
    mkdir -p "$(dirname "$WORKSPACE_CONF")"
    echo "ASH_WORKSPACE=$name" > "$WORKSPACE_CONF"
    ln -sfn "$WORKSPACE_DIR/$name" "$HOME/current-workspace"
    echo "Switched to workspace $name."
}

ash_up() {
    echo "Starting AI services..."
    sudo systemctl start qdrant || true
    sudo systemctl start ollama || true
    sudo systemctl start ash-lsfs || true
    echo "Services started."
}

ash_down() {
    echo "Stopping AI services..."
    sudo systemctl stop ash-lsfs || true
    sudo systemctl stop ollama || true
    sudo systemctl stop qdrant || true
    echo "Services stopped."
}

ash_status() {
    echo "AI Services Status:"
    sudo systemctl is-active qdrant && echo "Qdrant: Running" || echo "Qdrant: Stopped"
    sudo systemctl is-active ollama && echo "Ollama: Running" || echo "Ollama: Stopped"
    sudo systemctl is-active ash-lsfs && echo "LSFS: Running" || echo "LSFS: Stopped"
}

if [ "$1" = "workspace" ]; then
    shift
    case "$1" in
        create) create_workspace "$2" ;;
        reset) reset_workspace "$2" ;;
        list) list_workspaces ;;
        delete) delete_workspace "$2" ;;
        switch) switch_workspace "$2" ;;
        *) show_help ;;
    esac
elif [ "$1" = "up" ]; then
    ash_up
elif [ "$1" = "down" ]; then
    ash_down
elif [ "$1" = "status" ]; then
    ash_status
else
    show_help
fi

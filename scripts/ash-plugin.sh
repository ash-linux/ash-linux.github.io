#!/bin/bash
set -e

PLUGIN_DIR="/etc/ash/plugins"

cmd=$1
case "$cmd" in
    list)
        echo "Installed plugins:"
        for f in "$PLUGIN_DIR"/*.plugin; do
            if [ -f "$f" ]; then
                jq -r '"\(.name) [\(.trigger)] - \(.description)"' "$f"
            fi
        done
        ;;
    install)
        url=$2
        if [ -z "$url" ]; then echo "Usage: ash-plugin install <url>"; exit 1; fi
        wget -qO - "$url" | tar -xz -C "$PLUGIN_DIR"
        echo "Plugin installed."
        ;;
    remove)
        name=$2
        if [ -z "$name" ]; then echo "Usage: ash-plugin remove <name>"; exit 1; fi
        rm -f "$PLUGIN_DIR/$name.plugin" "$PLUGIN_DIR/$name.sh"
        echo "Plugin $name removed."
        ;;
    *)
        echo "Usage: ash-plugin list | install <url> | remove <name>"
        exit 1
        ;;
esac

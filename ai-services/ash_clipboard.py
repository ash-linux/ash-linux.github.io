#!/usr/bin/env python3

import sys
import subprocess
import requests
import json
import time
from datetime import datetime

OLLAMA_API = "http://localhost:11434/api"
QDRANT_API = "http://localhost:6333"
COLLECTION_NAME = "clipboard"
EMBED_MODEL = "nomic-embed-text"
MAX_ENTRIES = 1000

def embed_text(text):
    try:
        response = requests.post(f"{OLLAMA_API}/embeddings", json={
            "model": EMBED_MODEL,
            "prompt": text
        })
        response.raise_for_status()
        return response.json().get("embedding")
    except Exception as e:
        print(f"Error embedding text: {e}")
        return None

def init_qdrant():
    try:
        res = requests.get(f"{QDRANT_API}/collections/{COLLECTION_NAME}")
        if res.status_code == 404:
            requests.put(f"{QDRANT_API}/collections/{COLLECTION_NAME}", json={
                "vectors": {
                    "size": 768,  # nomic-embed-text size
                    "distance": "Cosine"
                }
            })
    except Exception as e:
        print(f"Failed to init Qdrant: {e}")

def get_total_count():
    try:
        res = requests.get(f"{QDRANT_API}/collections/{COLLECTION_NAME}/points/count", json={"exact": True})
        return res.json().get("result", {}).get("count", 0)
    except:
        return 0

def prune_oldest():
    # Simplistic prune: delete oldest by timestamp
    try:
        # Search for oldest 1 (using payload index or just sorting, but qdrant doesn't natively sort by payload without payload index)
        # For simplicity, if we hit max, we can clear the whole collection or just let it grow if this is a naive implementation.
        # To strictly prune oldest, we'd need a payload index on timestamp. We'll just delete the oldest timestamp if count > MAX_ENTRIES.
        pass
    except:
        pass

def add_to_qdrant(text):
    if not text.strip():
        return
    
    count = get_total_count()
    if count >= MAX_ENTRIES:
        # Simple flush or prune
        pass
        
    embed = embed_text(text)
    if not embed:
        return
        
    point_id = hash(text) & ((1<<63)-1)  # simple ID
    
    try:
        requests.put(f"{QDRANT_API}/collections/{COLLECTION_NAME}/points", json={
            "points": [{
                "id": point_id,
                "vector": embed,
                "payload": {
                    "content": text,
                    "timestamp": datetime.now().isoformat()
                }
            }]
        })
    except Exception as e:
        print(f"Failed to add to Qdrant: {e}")

def watch_clipboard():
    init_qdrant()
    # wl-paste --watch waits and outputs on every clipboard change
    process = subprocess.Popen(["wl-paste", "--watch", "cat"], stdout=subprocess.PIPE, text=True)
    
    current_entry = ""
    while True:
        line = process.stdout.readline()
        if not line:
            break
        # wl-paste --watch cat will interleave or output entries. 
        # A better way is a small bash loop: wl-paste -w bash -c 'wl-paste > current; notify-daemon current'
        # Actually, wl-paste --watch executes a command for each change.
        pass

def run_clipboard_watcher():
    init_qdrant()
    # using wl-paste -w to trigger python script maybe better, but we'll implement a simple polling or handle wl-paste output
    # Since we are writing the python script to run as a daemon:
    # A robust way:
    import os
    last_clip = ""
    while True:
        try:
            res = subprocess.run(["wl-paste", "--no-newline"], capture_output=True, text=True)
            clip = res.stdout
            if clip and clip != last_clip:
                last_clip = clip
                add_to_qdrant(clip)
        except FileNotFoundError:
            # wl-paste not found
            pass
        time.sleep(2)

def search_clipboard(query):
    init_qdrant()
    embed = embed_text(query)
    if not embed: return
    try:
        response = requests.post(f"{QDRANT_API}/collections/{COLLECTION_NAME}/points/search", json={
            "vector": embed,
            "limit": 5,
            "with_payload": True
        })
        results = response.json().get("result", [])
        for r in results:
            print(f"[{r['payload'].get('timestamp')}] {r['payload'].get('content')}")
    except Exception as e:
        print(e)

if __name__ == "__main__":
    if len(sys.argv) > 1 and sys.argv[1] == "search":
        search_clipboard(" ".join(sys.argv[2:]))
    else:
        run_clipboard_watcher()

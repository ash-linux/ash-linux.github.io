import os
import sys
import json
import requests
import argparse

def notify(message, level="info"):
    webhook_url = os.environ.get("ASH_WEBHOOK_URL")
    if not webhook_url:
        return
        
    payload = {"text": f"[{level.upper()}] Ash OS: {message}"}
    try:
        requests.post(webhook_url, json=payload, timeout=5)
    except Exception as e:
        pass

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("message")
    parser.add_argument("--level", default="info")
    args = parser.parse_args()
    notify(args.message, args.level)

#!/usr/bin/env python3

import sys
import argparse
import requests
import json
import os

OLLAMA_API = "http://localhost:11434/api"
QDRANT_API = "http://localhost:6333"
COLLECTION_NAME = "ash_fs"
EMBED_MODEL = "nomic-embed-text"
CHAT_MODEL = "llama3.2"

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
        sys.exit(1)

def search_qdrant(embedding, limit=10):
    try:
        response = requests.post(f"{QDRANT_API}/collections/{COLLECTION_NAME}/points/search", json={
            "vector": embedding,
            "limit": limit,
            "with_payload": True
        })
        response.raise_for_status()
        return response.json().get("result", [])
    except Exception as e:
        print(f"Error searching Qdrant: {e}")
        return []

def chat_with_context(question, context_chunks):
    context_text = "\n\n".join([f"Source: {c['payload'].get('filepath', 'Unknown')}\n{c['payload'].get('content', '')}" for c in context_chunks])
    
    prompt = f"""You are Ash, an AI assistant for the Ash Linux OS. Answer the question using the provided context. If the answer is not in the context, say so.
    
Context:
{context_text}

Question:
{question}
"""
    
    try:
        response = requests.post(f"{OLLAMA_API}/generate", json={
            "model": CHAT_MODEL,
            "prompt": prompt,
            "stream": False
        })
        response.raise_for_status()
        return response.json().get("response", "")
    except Exception as e:
        print(f"Error chatting with Ollama: {e}")
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(description="Ask Ash a question using RAG")
    parser.add_argument("question", help="The natural language question to ask")
    args = parser.parse_args()
    
    embedding = embed_text(args.question)
    results = search_qdrant(embedding)
    
    if not results:
        print("No relevant context found in Ash LSFS.")
        results = []
    
    answer = chat_with_context(args.question, results)
    
    print("\n=== Ash Answer ===")
    print(answer)
    print("\n=== Sources ===")
    seen = set()
    for res in results:
        filepath = res['payload'].get('filepath')
        if filepath and filepath not in seen:
            print(f"- {filepath}")
            seen.add(filepath)

if __name__ == "__main__":
    main()

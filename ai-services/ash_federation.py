import os
import sys
import yaml
import json
import asyncio
import aiohttp
from typing import List, Dict, Any

CONF_PATH = "/etc/ash/federation.conf"

def load_peers() -> List[Dict[str, Any]]:
    if not os.path.exists(CONF_PATH):
        return []
    try:
        with open(CONF_PATH, 'r') as f:
            data = yaml.safe_load(f)
            return data.get('peers', [])
    except Exception as e:
        print(f"Error loading federation config: {e}", file=sys.stderr)
        return []

async def fetch_peer(session: aiohttp.ClientSession, peer: Dict[str, Any], query_vector: List[float], limit: int) -> List[Dict[str, Any]]:
    host = peer.get('host')
    port = peer.get('port', 6333)
    api_key = peer.get('api_key')
    
    url = f"http://{host}:{port}/collections/lsfs/points/search"
    headers = {'Content-Type': 'application/json'}
    if api_key:
        headers['api-key'] = api_key
        
    payload = {
        "vector": query_vector,
        "limit": limit,
        "with_payload": True
    }
    
    try:
        async with session.post(url, headers=headers, json=payload, timeout=5.0) as response:
            if response.status == 200:
                data = await response.json()
                results = data.get('result', [])
                for r in results:
                    r['_source_machine'] = host
                return results
            else:
                return []
    except Exception as e:
        print(f"Error querying peer {host}: {e}", file=sys.stderr)
        return []

def rrf(results_list: List[List[Dict[str, Any]]], k=60) -> List[Dict[str, Any]]:
    scores = {}
    item_map = {}
    
    for results in results_list:
        for rank, item in enumerate(results):
            doc_id = item.get('id')
            host = item.get('_source_machine')
            key = f"{host}_{doc_id}"
            
            if key not in scores:
                scores[key] = 0
                item_map[key] = item
                
            scores[key] += 1.0 / (k + rank + 1)
            
    ranked = sorted(scores.items(), key=lambda x: x[1], reverse=True)
    
    final_results = []
    for key, score in ranked:
        item = item_map[key]
        item['_rrf_score'] = score
        final_results.append(item)
        
    return final_results

async def federated_search(query_vector: List[float], limit: int = 10) -> List[Dict[str, Any]]:
    peers = load_peers()
    if not peers:
        return []
        
    async with aiohttp.ClientSession() as session:
        tasks = [fetch_peer(session, peer, query_vector, limit) for peer in peers]
        results_list = await asyncio.gather(*tasks)
        
    merged = rrf(results_list)
    return merged[:limit]

def main():
    if len(sys.argv) < 2:
        print("Usage: ash_federation.py <json_query_vector_file>")
        sys.exit(1)
        
    vector_file = sys.argv[1]
    with open(vector_file, 'r') as f:
        vector = json.load(f)
        
    limit = 10
    if len(sys.argv) > 2:
        limit = int(sys.argv[2])
        
    results = asyncio.run(federated_search(vector, limit))
    print(json.dumps(results, indent=2))

if __name__ == '__main__':
    main()

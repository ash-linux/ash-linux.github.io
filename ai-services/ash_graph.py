import os
import sqlite3
import sys
import re
import argparse

DB_PATH = os.path.expanduser("~/.local/share/ash/graph.db")

def init_db():
    os.makedirs(os.path.dirname(DB_PATH), exist_ok=True)
    conn = sqlite3.connect(DB_PATH)
    c = conn.cursor()
    c.execute("""
        CREATE TABLE IF NOT EXISTS edges (
            source TEXT,
            target TEXT,
            rel_type TEXT,
            UNIQUE(source, target, rel_type)
        )
    """)
    conn.commit()
    return conn

def add_edge(source, target, rel_type):
    conn = init_db()
    c = conn.cursor()
    c.execute("INSERT OR IGNORE INTO edges (source, target, rel_type) VALUES (?, ?, ?)", (source, target, rel_type))
    conn.commit()
    conn.close()

def parse_imports(filepath):
    # simplistic import parsing for Python and JS
    if not os.path.exists(filepath): return
    ext = os.path.splitext(filepath)[1]
    
    try:
        with open(filepath, 'r') as f:
            content = f.read()
    except Exception:
        return
        
    if ext == '.py':
        for match in re.finditer(r'^(?:from|import)\s+([a-zA-Z0-9_\.]+)', content, re.MULTILINE):
            target = match.group(1)
            add_edge(filepath, target, 'imports')
    elif ext in ['.js', '.ts', '.jsx', '.tsx']:
        for match in re.finditer(r'import\s+.*?\s+from\s+[\'"](.*?)[\'"]', content, re.MULTILINE):
            target = match.group(1)
            add_edge(filepath, target, 'imports')
        for match in re.finditer(r'require\([\'"](.*?)[\'"]\)', content):
            target = match.group(1)
            add_edge(filepath, target, 'imports')

def build_git_coedits(repo_path):
    # This would use git log --name-only to find files edited in the same commit
    pass

def query_graph(file, direction='both'):
    conn = init_db()
    c = conn.cursor()
    results = []
    
    if direction in ['out', 'both']:
        c.execute("SELECT target, rel_type FROM edges WHERE source = ?", (file,))
        for row in c.fetchall():
            results.append(f"-> {row[0]} ({row[1]})")
            
    if direction in ['in', 'both']:
        c.execute("SELECT source, rel_type FROM edges WHERE target = ?", (file,))
        for row in c.fetchall():
            results.append(f"<- {row[0]} ({row[1]})")
            
    conn.close()
    return results

def main():
    parser = argparse.ArgumentParser()
    subparsers = parser.add_subparsers(dest='cmd')
    
    scan_p = subparsers.add_parser('scan')
    scan_p.add_argument('file')
    
    query_p = subparsers.add_parser('query')
    query_p.add_argument('file')
    
    deps_p = subparsers.add_parser('deps')
    deps_p.add_argument('file')
    
    args = parser.parse_args()
    
    if args.cmd == 'scan':
        parse_imports(args.file)
    elif args.cmd == 'query':
        res = query_graph(args.file)
        print(f"Relationships for {args.file}:")
        for r in res: print(r)
    elif args.cmd == 'deps':
        res = query_graph(args.file, direction='out')
        print(f"Dependencies of {args.file}:")
        for r in res: print(r)

if __name__ == '__main__':
    main()

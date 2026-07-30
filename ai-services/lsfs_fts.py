import sqlite3
import os
import time

DB_PATH = os.environ.get("LSFS_FTS_DB", os.path.expanduser("~/.local/share/ash/lsfs_fts.db"))

def init_db():
    os.makedirs(os.path.dirname(DB_PATH), exist_ok=True)
    conn = sqlite3.connect(DB_PATH)
    conn.execute("""
        CREATE VIRTUAL TABLE IF NOT EXISTS fts_index USING fts5(
            path UNINDEXED,
            name,
            content,
            tokenize='porter'
        )
    """)
    conn.execute("CREATE TABLE IF NOT EXISTS sync_state (path TEXT PRIMARY KEY, mtime INTEGER)")
    conn.commit()
    return conn

def update_index(path, name, content, mtime):
    conn = init_db()
    c = conn.cursor()
    c.execute("DELETE FROM fts_index WHERE path = ?", (path,))
    c.execute("INSERT INTO fts_index (path, name, content) VALUES (?, ?, ?)", (path, name, content))
    c.execute("INSERT OR REPLACE INTO sync_state (path, mtime) VALUES (?, ?)", (path, mtime))
    conn.commit()
    conn.close()

def delete_index(path):
    conn = init_db()
    c = conn.cursor()
    c.execute("DELETE FROM fts_index WHERE path = ?", (path,))
    c.execute("DELETE FROM sync_state WHERE path = ?", (path,))
    conn.commit()
    conn.close()

def bm25_search(query, limit=20):
    conn = init_db()
    c = conn.cursor()
    try:
        # FTS5 requires quotes or similar for some queries, we'll strip unneeded chars and just pass it
        clean_query = query.replace('"', '""')
        c.execute("""
            SELECT path, name, bm25(fts_index) as score
            FROM fts_index
            WHERE fts_index MATCH ?
            ORDER BY score
            LIMIT ?
        """, (f'"{clean_query}"*', limit))
        # score from bm25 is negative by default (lower is better in SQLite FTS5)
        # we will return absolute or inverted score
        results = [{"path": row[0], "name": row[1], "score": abs(row[2])} for row in c.fetchall()]
        return results
    except Exception as e:
        return []
    finally:
        conn.close()

if __name__ == "__main__":
    import sys
    if len(sys.argv) > 1:
        query = " ".join(sys.argv[1:])
        results = bm25_search(query)
        for r in results:
            print(f"{r['path']} | {r['score']}")

CREATE TABLE IF NOT EXISTS files (
    id TEXT PRIMARY KEY,
    path TEXT NOT NULL,
    file_hash TEXT NOT NULL UNIQUE,
    file_name TEXT NOT NULL,
    description TEXT,
    size INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS temporary_links (
    id TEXT PRIMARY KEY,
    file_id TEXT NOT NULL,
    token TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (file_id) REFERENCES files (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS file_stats (
    id TEXT PRIMARY KEY,
    file_id TEXT NOT NULL UNIQUE,
    download_count INTEGER DEFAULT 0,
    last_downloaded_at DATETIME,
    FOREIGN KEY (file_id) REFERENCES files (id) ON DELETE CASCADE
);

CREATE INDEX idx_temporary_links_file_id ON temporary_links (file_id);

CREATE INDEX idx_file_stats_file_id ON file_stats (file_id);

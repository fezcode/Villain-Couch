CREATE TABLE IF NOT EXISTS media_files (
    "filepath" TEXT NOT NULL PRIMARY KEY,
    "filename" TEXT,
    "total_seconds" INTEGER,
    "current_second" INTEGER,
    "created_at" DATETIME,
    "updated_at" DATETIME
);

-- Create the main 'workspaces' table.
CREATE TABLE IF NOT EXISTS workspaces (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    directory_path TEXT NOT NULL UNIQUE,
    directory_name TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now')),
    updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now')),
    deleted_at TEXT DEFAULT NULL
);
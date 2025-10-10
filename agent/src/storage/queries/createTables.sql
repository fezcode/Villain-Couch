CREATE TABLE IF NOT EXISTS media_files (
    "filepath" TEXT NOT NULL PRIMARY KEY,
    "filename" TEXT,
    "total_seconds" INTEGER,
    "current_second" INTEGER,
    "created_at" DATETIME,
    "updated_at" DATETIME
);
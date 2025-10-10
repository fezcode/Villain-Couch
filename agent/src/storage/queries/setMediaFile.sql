-- This query handles the UPSERT logic.
-- For INSERT: We provide all values. created_at will be set to the current time.
-- For UPDATE (ON CONFLICT): We only update the columns that should change, leaving the existing created_at value untouched.
INSERT INTO media_files (filepath, filename, total_seconds, current_second, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?)
    ON CONFLICT(filepath) DO UPDATE SET
    filename = excluded.filename,
    total_seconds = excluded.total_seconds,
    current_second = excluded.current_second,
    updated_at = excluded.updated_at;
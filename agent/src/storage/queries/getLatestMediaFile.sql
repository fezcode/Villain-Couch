SELECT filepath, filename, total_seconds, current_second, created_at, updated_at
    FROM media_files
    ORDER BY updated_at DESC
    LIMIT 1;

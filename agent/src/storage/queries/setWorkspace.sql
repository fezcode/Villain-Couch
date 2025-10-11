INSERT INTO workspaces (directory_path, directory_name, created_at, updated_at)
VALUES (?, ?, ?, ?) ON CONFLICT(directory_path) DO NOTHING
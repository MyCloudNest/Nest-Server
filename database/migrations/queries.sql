-- name: InsertFile :one
INSERT INTO
    files (id, file_name, file_hash, path, description, size)
VALUES
    (?, ?, ?, ?, ?, ?) RETURNING id,
    created_at;

-- name: GetFileByID :one
SELECT
    id,
    file_name,
    path,
    description,
    size,
    created_at
FROM
    files
WHERE
    id = ?;

-- name: GetFileByHash :one
SELECT
    id,
    path,
    file_name
FROM
    files
WHERE
    file_hash = ?;

-- name: GetAllFiles :many
SELECT
    id,
    path,
    file_name,
    description,
    size,
    created_at
FROM
    files
ORDER BY
    created_at DESC;

-- name: UpdateFileDescription :one
UPDATE files
SET
    description = ?
WHERE
    id = ? RETURNING id,
    description;

-- name: DeleteFile :one
DELETE FROM files
WHERE
    id = ? RETURNING id;

-- name: InsertTemporaryLink :one
INSERT INTO
    temporary_links (id, file_id, token, expires_at)
VALUES
    (?, ?, ?, ?) RETURNING id,
    token,
    created_at,
    expires_at;

-- name: GetTemporaryLinkByToken :one
SELECT
    id,
    file_id,
    token,
    created_at,
    expires_at
FROM
    temporary_links
WHERE
    token = ?;

-- name: DeleteExpiredLinks :exec
DELETE FROM temporary_links
WHERE
    expires_at < CURRENT_TIMESTAMP;

-- name: UpsertFileStats :one
INSERT INTO
    file_stats (id, file_id, download_count, last_downloaded_at)
VALUES
    (?, ?, ?, ?) ON CONFLICT (file_id) DO
UPDATE
SET
    download_count = download_count + excluded.download_count,
    last_downloaded_at = excluded.last_downloaded_at RETURNING id,
    download_count,
    last_downloaded_at;

-- name: GetFileStats :one
SELECT
    id,
    file_id,
    download_count,
    last_downloaded_at
FROM
    file_stats
WHERE
    file_id = ?;

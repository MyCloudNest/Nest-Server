// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const deleteExpiredLinks = `-- name: DeleteExpiredLinks :exec
DELETE FROM temporary_links
WHERE
    expires_at < CURRENT_TIMESTAMP
`

func (q *Queries) DeleteExpiredLinks(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteExpiredLinks)
	return err
}

const deleteFile = `-- name: DeleteFile :one
DELETE FROM files
WHERE
    id = ? RETURNING id
`

func (q *Queries) DeleteFile(ctx context.Context, id string) (string, error) {
	row := q.db.QueryRowContext(ctx, deleteFile, id)
	err := row.Scan(&id)
	return id, err
}

const deleteFilePermissions = `-- name: DeleteFilePermissions :exec
DELETE FROM permissions
WHERE
    file_id = ?
`

func (q *Queries) DeleteFilePermissions(ctx context.Context, fileID string) error {
	_, err := q.db.ExecContext(ctx, deleteFilePermissions, fileID)
	return err
}

const getAllFiles = `-- name: GetAllFiles :many
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
    created_at DESC
`

type GetAllFilesRow struct {
	ID          string
	Path        string
	FileName    string
	Description sql.NullString
	Size        int64
	CreatedAt   sql.NullTime
}

func (q *Queries) GetAllFiles(ctx context.Context) ([]GetAllFilesRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllFiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllFilesRow
	for rows.Next() {
		var i GetAllFilesRow
		if err := rows.Scan(
			&i.ID,
			&i.Path,
			&i.FileName,
			&i.Description,
			&i.Size,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFileByHash = `-- name: GetFileByHash :one
SELECT
    id,
    path,
    file_name
FROM
    files
WHERE
    file_hash = ?
`

type GetFileByHashRow struct {
	ID       string
	Path     string
	FileName string
}

func (q *Queries) GetFileByHash(ctx context.Context, fileHash string) (GetFileByHashRow, error) {
	row := q.db.QueryRowContext(ctx, getFileByHash, fileHash)
	var i GetFileByHashRow
	err := row.Scan(&i.ID, &i.Path, &i.FileName)
	return i, err
}

const getFileByID = `-- name: GetFileByID :one
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
    id = ?
`

type GetFileByIDRow struct {
	ID          string
	FileName    string
	Path        string
	Description sql.NullString
	Size        int64
	CreatedAt   sql.NullTime
}

func (q *Queries) GetFileByID(ctx context.Context, id string) (GetFileByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getFileByID, id)
	var i GetFileByIDRow
	err := row.Scan(
		&i.ID,
		&i.FileName,
		&i.Path,
		&i.Description,
		&i.Size,
		&i.CreatedAt,
	)
	return i, err
}

const getFilePermissions = `-- name: GetFilePermissions :many
SELECT
    id,
    file_id,
    permission_type,
    granted_at
FROM
    permissions
WHERE
    file_id = ?
`

func (q *Queries) GetFilePermissions(ctx context.Context, fileID string) ([]Permission, error) {
	rows, err := q.db.QueryContext(ctx, getFilePermissions, fileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Permission
	for rows.Next() {
		var i Permission
		if err := rows.Scan(
			&i.ID,
			&i.FileID,
			&i.PermissionType,
			&i.GrantedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFileStats = `-- name: GetFileStats :one
SELECT
    id,
    file_id,
    download_count,
    last_downloaded_at
FROM
    file_stats
WHERE
    file_id = ?
`

func (q *Queries) GetFileStats(ctx context.Context, fileID string) (FileStat, error) {
	row := q.db.QueryRowContext(ctx, getFileStats, fileID)
	var i FileStat
	err := row.Scan(
		&i.ID,
		&i.FileID,
		&i.DownloadCount,
		&i.LastDownloadedAt,
	)
	return i, err
}

const getTemporaryLinkByToken = `-- name: GetTemporaryLinkByToken :one
SELECT
    id,
    file_id,
    token,
    created_at,
    expires_at
FROM
    temporary_links
WHERE
    token = ?
`

func (q *Queries) GetTemporaryLinkByToken(ctx context.Context, token string) (TemporaryLink, error) {
	row := q.db.QueryRowContext(ctx, getTemporaryLinkByToken, token)
	var i TemporaryLink
	err := row.Scan(
		&i.ID,
		&i.FileID,
		&i.Token,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const insertFile = `-- name: InsertFile :one
INSERT INTO
    files (id, file_name, file_hash, path, description, size)
VALUES
    (?, ?, ?, ?, ?, ?) RETURNING id,
    created_at
`

type InsertFileParams struct {
	ID          string
	FileName    string
	FileHash    string
	Path        string
	Description sql.NullString
	Size        int64
}

type InsertFileRow struct {
	ID        string
	CreatedAt sql.NullTime
}

func (q *Queries) InsertFile(ctx context.Context, arg InsertFileParams) (InsertFileRow, error) {
	row := q.db.QueryRowContext(ctx, insertFile,
		arg.ID,
		arg.FileName,
		arg.FileHash,
		arg.Path,
		arg.Description,
		arg.Size,
	)
	var i InsertFileRow
	err := row.Scan(&i.ID, &i.CreatedAt)
	return i, err
}

const insertPermission = `-- name: InsertPermission :one
INSERT INTO
    permissions (id, file_id, permission_type)
VALUES
    (?, ?, ?) RETURNING id,
    granted_at
`

type InsertPermissionParams struct {
	ID             string
	FileID         string
	PermissionType string
}

type InsertPermissionRow struct {
	ID        string
	GrantedAt sql.NullTime
}

func (q *Queries) InsertPermission(ctx context.Context, arg InsertPermissionParams) (InsertPermissionRow, error) {
	row := q.db.QueryRowContext(ctx, insertPermission, arg.ID, arg.FileID, arg.PermissionType)
	var i InsertPermissionRow
	err := row.Scan(&i.ID, &i.GrantedAt)
	return i, err
}

const insertTemporaryLink = `-- name: InsertTemporaryLink :one
INSERT INTO
    temporary_links (id, file_id, token, expires_at)
VALUES
    (?, ?, ?, ?) RETURNING id,
    token,
    created_at,
    expires_at
`

type InsertTemporaryLinkParams struct {
	ID        string
	FileID    string
	Token     string
	ExpiresAt time.Time
}

type InsertTemporaryLinkRow struct {
	ID        string
	Token     string
	CreatedAt sql.NullTime
	ExpiresAt time.Time
}

func (q *Queries) InsertTemporaryLink(ctx context.Context, arg InsertTemporaryLinkParams) (InsertTemporaryLinkRow, error) {
	row := q.db.QueryRowContext(ctx, insertTemporaryLink,
		arg.ID,
		arg.FileID,
		arg.Token,
		arg.ExpiresAt,
	)
	var i InsertTemporaryLinkRow
	err := row.Scan(
		&i.ID,
		&i.Token,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const updateFileDescription = `-- name: UpdateFileDescription :one
UPDATE files
SET
    description = ?
WHERE
    id = ? RETURNING id,
    description
`

type UpdateFileDescriptionParams struct {
	Description sql.NullString
	ID          string
}

type UpdateFileDescriptionRow struct {
	ID          string
	Description sql.NullString
}

func (q *Queries) UpdateFileDescription(ctx context.Context, arg UpdateFileDescriptionParams) (UpdateFileDescriptionRow, error) {
	row := q.db.QueryRowContext(ctx, updateFileDescription, arg.Description, arg.ID)
	var i UpdateFileDescriptionRow
	err := row.Scan(&i.ID, &i.Description)
	return i, err
}

const upsertFileStats = `-- name: UpsertFileStats :one
INSERT INTO
    file_stats (id, file_id, download_count, last_downloaded_at)
VALUES
    (?, ?, ?, ?) ON CONFLICT (file_id) DO
UPDATE
SET
    download_count = download_count + excluded.download_count,
    last_downloaded_at = excluded.last_downloaded_at RETURNING id,
    download_count,
    last_downloaded_at
`

type UpsertFileStatsParams struct {
	ID               string
	FileID           string
	DownloadCount    sql.NullInt64
	LastDownloadedAt sql.NullTime
}

type UpsertFileStatsRow struct {
	ID               string
	DownloadCount    sql.NullInt64
	LastDownloadedAt sql.NullTime
}

func (q *Queries) UpsertFileStats(ctx context.Context, arg UpsertFileStatsParams) (UpsertFileStatsRow, error) {
	row := q.db.QueryRowContext(ctx, upsertFileStats,
		arg.ID,
		arg.FileID,
		arg.DownloadCount,
		arg.LastDownloadedAt,
	)
	var i UpsertFileStatsRow
	err := row.Scan(&i.ID, &i.DownloadCount, &i.LastDownloadedAt)
	return i, err
}
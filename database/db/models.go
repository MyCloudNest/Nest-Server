// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"time"
)

type File struct {
	ID          string
	Path        string
	FileHash    string
	FileName    string
	Description sql.NullString
	Size        int64
	CreatedAt   sql.NullTime
}

type FileStat struct {
	ID               string
	FileID           string
	DownloadCount    sql.NullInt64
	LastDownloadedAt sql.NullTime
}

type Permission struct {
	ID             string
	FileID         string
	PermissionType string
	GrantedAt      sql.NullTime
}

type TemporaryLink struct {
	ID        string
	FileID    string
	Token     string
	CreatedAt sql.NullTime
	ExpiresAt time.Time
}

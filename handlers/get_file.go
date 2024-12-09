package handlers

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/MyCloudNest/Nest-Server/database"
	"github.com/MyCloudNest/Nest-Server/database/db"
	"github.com/gofiber/fiber/v2"
)

func GetFile(c *fiber.Ctx) error {
	fileID := c.Params("id")
	download := c.Query("download", "false")

	if fileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"ok":      false,
			"message": "File ID is required",
		})
	}

	queries := db.New(database.DB)

	file, err := queries.GetFileByID(c.Context(), fileID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"ok":      false,
			"message": "File not found",
		})
	}

	if _, err = os.Stat(file.Path); os.IsNotExist(err) {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"ok":      false,
			"message": "File not found on disk",
		})
	}

	statsID := fmt.Sprintf("%s-stats", fileID)
	_, err = queries.UpsertFileStats(c.Context(), db.UpsertFileStatsParams{
		ID:               statsID,
		FileID:           fileID,
		DownloadCount:    sql.NullInt64{Int64: 1, Valid: true},
		LastDownloadedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": "Failed to update file stats",
		})
	}

	if download == "true" {
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileName))
	} else {
		c.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.FileName))
	}

	return c.SendFile(file.Path)
}

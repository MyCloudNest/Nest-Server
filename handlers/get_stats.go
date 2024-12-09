package handlers

import (
	"github.com/MyCloudNest/Nest-Server/database"
	"github.com/MyCloudNest/Nest-Server/database/db"
	"github.com/gofiber/fiber/v2"
)

func GetStats(c *fiber.Ctx) error {
	fileId := c.Params("id")

	if fileId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"ok":      false,
			"message": "File ID is required",
		})
	}

	queries := db.New(database.DB)

	stats, err := queries.GetFileStats(c.Context(), fileId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"ok":      false,
			"message": "File not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"ok":                 true,
		"file_id":            stats.FileID,
		"download_count":     stats.DownloadCount,
		"last_downloaded_at": stats.LastDownloadedAt,
	})
}

package handlers

import (
	"os"

	"github.com/MyCloudNest/Nest-Server/database"
	"github.com/MyCloudNest/Nest-Server/database/db"
	"github.com/gofiber/fiber/v2"
)

func DeleteFile(c *fiber.Ctx) error {
	fileID := c.Params("id")
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

	if err := os.Remove(file.Path); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": "Failed to delete file from disk",
		})
	}

	if _, err := queries.DeleteFile(c.Context(), fileID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": "Failed to delete file record from database",
		})
	}

	return c.JSON(&fiber.Map{
		"ok":      true,
		"message": "File deleted successfully",
	})
}

package handlers

import (
	"github.com/MyCloudNest/Nest-Server/database"
	"github.com/MyCloudNest/Nest-Server/database/db"
	"github.com/MyCloudNest/Nest-Server/schemas/requests"
	"github.com/gofiber/fiber/v2"
)

func RetrieveFiles(c *fiber.Ctx) error {
	queries := db.New(database.DB)

	files, err := queries.GetAllFiles(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": "Failed to retrieve files",
		})
	}

	var response []requests.FileResponse
	for _, file := range files {
		response = append(response, requests.FileResponse{
			ID:        file.ID,
			FileName:  file.FileName,
			Path:      file.Path,
			Size:      file.Size,
			CreatedAt: file.CreatedAt.Time.Format("2006-01-02T15:04:05Z"),
		})
	}

	return c.JSON(&fiber.Map{
		"ok":    true,
		"files": response,
	})
}

package handlers

import (
	"fmt"
	"time"

	"github.com/MyCloudNest/Nest-Server/database"
	"github.com/MyCloudNest/Nest-Server/database/db"
	"github.com/gofiber/fiber/v2"
)

func ValidateTempLink(c *fiber.Ctx) error {
	token := c.Query("token")
	download := c.Query("download", "false")

	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"ok":      false,
			"message": "Token is required",
		})
	}

	queries := db.New(database.DB)

	tempLink, err := queries.GetTemporaryLinkByToken(c.Context(), token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"ok":      false,
			"message": "Invalid or expired token",
		})
	}

	if time.Now().After(tempLink.ExpiresAt) {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"ok":      false,
			"message": "Link expired",
		})
	}

	file, err := queries.GetFileByID(c.Context(), tempLink.FileID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"ok":      false,
			"message": "File not found",
		})
	}

	if download == "true" {
		c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.FileName))
	} else {
		c.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", file.FileName))
	}

	return c.SendFile(file.Path)
}

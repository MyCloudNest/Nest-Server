package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/MyCloudNest/Nest-Server/database"
	"github.com/MyCloudNest/Nest-Server/database/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GenerateTempLink(c *fiber.Ctx) error {
	fileID := c.Params("id")
	if fileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"ok":      false,
			"message": "File ID is required",
		})
	}

	expirationStr := c.Query("expires_at", "3600")
	if expirationStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"ok":      false,
			"message": "Expiration timestamp is required",
		})
	}

	unixTimestamp, err := strconv.ParseInt(expirationStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"ok":      false,
			"message": "Invalid expiration timestamp format. Use Unix timestamp format",
		})
	}

	currentTime := time.Now().Unix()
	expiration := time.Unix(currentTime+unixTimestamp, 0)

	queries := db.New(database.DB)

	_, err = queries.GetFileByID(c.Context(), fileID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"ok":      false,
			"message": "File not found",
		})
	}

	token := generateToken()

	_, err = queries.InsertTemporaryLink(c.Context(), db.InsertTemporaryLinkParams{
		ID:        uuid.New().String(),
		Token:     token,
		FileID:    fileID,
		ExpiresAt: expiration,
	})
	if err != nil {
		fmt.Printf("Error creating temporary link: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": fmt.Sprintf("Failed to create temporary link: %v", err),
		})
	}

	tempLink := fmt.Sprintf("http://localhost:3000/api/v1/files/download?token=%s", token)

	return c.JSON(&fiber.Map{
		"ok":         true,
		"link":       tempLink,
		"expires_at": expiration.Unix(),
	})
}

func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

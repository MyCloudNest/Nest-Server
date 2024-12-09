package handlers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/MyCloudNest/Nest-Server/database"
	"github.com/MyCloudNest/Nest-Server/database/db"
	"github.com/MyCloudNest/Nest-Server/schemas/requests"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UploadFile(c *fiber.Ctx) error {
	var req requests.UploadFileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"ok":      false,
			"message": "Invalid request body",
		})
	}

	id := uuid.New().String()[:8]

	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"ok":      false,
			"message": "Invalid request body",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"ok":      false,
			"message": "File is required",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": "Failed to open file",
		})
	}
	defer src.Close()

	hash := sha256.New()
	if _, err = io.Copy(hash, src); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": "Failed to calculate file hash",
		})
	}
	fileHash := hex.EncodeToString(hash.Sum(nil))

	queries := db.New(database.DB)
	existingFile, err := queries.GetFileByHash(c.Context(), fileHash)
	if err == nil {
		return c.JSON(fiber.Map{
			"ok":        true,
			"message":   "File already exists",
			"file_id":   existingFile.ID,
			"file_name": existingFile.FileName,
			"file_path": existingFile.Path,
		})
	}

	if req.FileName == "" {
		req.FileName = file.Filename
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": "Error getting home directory",
		})
	}

	baseDir := filepath.Join(homeDir, ".cloudnest")
	var targetDir string

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	switch fileExt {
	case ".mp3", ".wav", ".flac", ".ogg":
		targetDir = filepath.Join(baseDir, "audio")
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		targetDir = filepath.Join(baseDir, "image")
	case ".mp4", ".avi", ".mov", ".mkv":
		targetDir = filepath.Join(baseDir, "video")
	case ".pdf", ".doc", ".docx", ".txt", ".xlsx", ".go", ".py", ".js", ".html", ".css", ".java", ".c", ".cpp", ".ts", ".json", ".xml", ".sql", ".sh", ".rs", ".bat", ".zip", ".rar":
		targetDir = filepath.Join(baseDir, "document")
	default:
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"ok":      false,
			"message": "Unsupported file type",
		})
	}

	fullPath := filepath.Join(targetDir, req.FileName)

	if _, err = os.Stat(fullPath); err == nil {
		return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
			"ok":      false,
			"message": "File already exists",
		})
	}

	if fullPath != "" {
		if err = os.MkdirAll(filepath.Dir(fullPath), 0o700); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"ok":      false,
				"message": "Failed to create directory",
			})
		}
	}

	newFileName := id + "_" + req.FileName
	fullPath = filepath.Join(filepath.Dir(fullPath), newFileName)

	if err = c.SaveFile(file, fullPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": "Failed to save file",
		})
	}

	log.Printf("File uploaded: %s to %s", file.Filename, fullPath)

	_, err = queries.InsertFile(c.Context(), db.InsertFileParams{
		ID:          id,
		Path:        fullPath,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
		Size:        file.Size,
		FileHash:    fileHash,
		FileName:    req.FileName,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"ok":      false,
			"message": "Failed to save file metadata",
		})
	}

	return c.JSON(fiber.Map{
		"ok":        true,
		"file_id":   id,
		"file_name": req.FileName,
		"file_path": fullPath,
		"file_hash": fileHash,
	})
}

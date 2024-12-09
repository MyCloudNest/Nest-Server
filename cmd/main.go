package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/MyCloudNest/Nest-Server/config"
	"github.com/MyCloudNest/Nest-Server/database"
	"github.com/MyCloudNest/Nest-Server/handlers"
	"github.com/MyCloudNest/Nest-Server/middleware"
	"github.com/MyCloudNest/Nest-Server/utils"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/redis"
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	cloudNestDir := filepath.Join(homeDir, ".cloudnest")
	logDir := filepath.Join(cloudNestDir, "log")
	configPath := filepath.Join(cloudNestDir, "config.toml")
	subDirs := []string{"audio", "image", "video", "document"}

	if err = os.MkdirAll(cloudNestDir, 0o700); err != nil {
		log.Fatalf("Failed to create base directory: %v", err)
	}

	for _, dir := range subDirs {
		dirPath := filepath.Join(cloudNestDir, dir)
		if err = os.MkdirAll(dirPath, 0o700); err != nil {
			log.Fatalf("Failed to create %s directory: %v", dir, err)
		}
	}

	if err = os.MkdirAll(logDir, 0o700); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	errorLogPath := filepath.Join(logDir, "error.log")
	outputLogPath := filepath.Join(logDir, "output.log")

	errorLogFile, err := os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		log.Fatalf("Failed to create error log file: %v", err)
	}

	outputLogFile, err := os.OpenFile(outputLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		log.Fatalf("Failed to create output log file: %v", err)
	}

	log.SetOutput(outputLogFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	os.Stderr = errorLogFile

	configURL := "https://raw.githubusercontent.com/MyCloudNest/Nest-Server/refs/heads/main/config/config.toml"
	if err = utils.DownloadFile(configURL, configPath); err != nil {
		log.Fatalf("Failed to download config file: %v", err)
	}

	log.Println("Initialization complete. Config file downloaded to:", configPath)
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		Prefork:     cfg.Performance.PerFork,
		BodyLimit:   cfg.RateLimit.LimitBody,
		Concurrency: cfg.Performance.Concurrency,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"ok":      false,
				"message": "Internal server error",
			})
		},
	})

	if cfg.Whitelist.Enabled {
		app.Use(middleware.WhitelistMiddleware(cfg.Whitelist.WhitelistedIPs))
	}

	app.Use(compress.New())
	app.Use(logger.New())
	app.Use(recover.New())

	if cfg.RateLimit.Enabled {

		redisStore := redis.New(redis.Config{
			Host:      "localhost",
			Port:      6379,
			Username:  "",
			Password:  "",
			Database:  0,
			Reset:     false,
			TLSConfig: nil,
			PoolSize:  10,
		})

		defer redisStore.Close()

		app.Use(limiter.New(limiter.Config{
			Max:        cfg.RateLimit.MaxRequests,
			Expiration: time.Duration(cfg.RateLimit.ExpireTime),
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(429).JSON(fiber.Map{
					"ok":      false,
					"message": "Too many requests",
				})
			},
			Storage: redisStore,
		}))
	}

	if cfg.Cache.Enabled {
		app.Use(cache.New(cache.Config{
			Next: func(c *fiber.Ctx) bool {
				return c.Query("refresh") == "true"
			},
			Expiration:           10 * time.Minute, // 10 minutes
			CacheControl:         true,
			StoreResponseHeaders: true,
			Storage:              nil,
			MaxBytes:             100 * 1024 * 1024, // 100MB
			ExpirationGenerator: func(c *fiber.Ctx, cfg *cache.Config) time.Duration {
				if c.Path() == "/api/v1/files" && c.Query("file_id") != "" {
					return 30 * time.Minute // 30 minutes
				}
				return cfg.Expiration
			},
			CacheHeader: "X-Cache",
		}))
	}

	// close database safely on exit
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-exitSignal
		log.Println("Shutting down server...")
		database.CloseDB()
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error shutting down server: %v", err)
		}

		log.Println("Server stopped safely")
		os.Exit(0)
	}()

	app.Post("/api/v1/files", handlers.UploadFile)

	app.Get("/api/v1/files/download", handlers.ValidateTempLink)

	app.Get("/api/v1/files/:id", handlers.GetFile)

	app.Get("/api/v1/files", handlers.RetrieveFiles)

	app.Delete("/api/v1/files/:id", handlers.DeleteFile)

	app.Post("/api/v1/files/:id/temp-link", handlers.GenerateTempLink)

	app.Get("/api/v1/files/:id/stats", handlers.GetStats)

	log.Fatal(app.Listen(cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)))
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/sqlite3"
	"github.com/run-llama/study-llama/frontend/auth"
	"github.com/run-llama/study-llama/frontend/handlers"
)

func main() {
	// Create a new Fiber app
	app := Setup()

	// Start the Fiber server on port 8000
	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func cacheSetupPost(keyGen func(*fiber.Ctx) string) fiber.Handler {
	cacheStorage := sqlite3.New(
		sqlite3.Config{
			Database:        "cache_post.db",
			Table:           os.Getenv("CACHE_TABLE"),
			ConnMaxLifetime: 5 * time.Second,
		},
	)
	cache := cache.New(
		cache.Config{
			Expiration:   1 * time.Hour,
			CacheControl: true,
			Storage:      cacheStorage,
			KeyGenerator: keyGen,
			Methods:      []string{fiber.MethodPost},
		},
	)
	return cache
}

func cacheSetupGet(keyGen func(*fiber.Ctx) string) fiber.Handler {
	cacheStorage := sqlite3.New(
		sqlite3.Config{
			Database:        "cache_get.db",
			Table:           os.Getenv("CACHE_TABLE"),
			ConnMaxLifetime: 5 * time.Second,
		},
	)
	cache := cache.New(
		cache.Config{
			Expiration:   24 * time.Hour,
			CacheControl: true,
			Storage:      cacheStorage,
			KeyGenerator: keyGen,
			Methods:      []string{fiber.MethodGet, fiber.MethodHead},
		},
	)
	return cache
}

func corsSetup(methods string) fiber.Handler {
	allowedOrigins := []string{"https://gityear.re"}
	corsHandler := cors.New(
		cors.Config{
			AllowOriginsFunc: func(origin string) bool {
				return slices.Contains(allowedOrigins, origin)
			},
			AllowMethods: methods,
		},
	)
	return corsHandler
}

func limiterSetup(reqPerMinute int) fiber.Handler {
	limiterStorage := sqlite3.New(
		sqlite3.Config{
			Database:        "ratelimiter.db",
			Table:           os.Getenv("RATE_LIMITING_TABLE"),
			ConnMaxLifetime: 5 * time.Second,
		},
	)
	limiter := limiter.New(
		limiter.Config{
			Max:     reqPerMinute,
			Storage: limiterStorage,
		},
	)
	return limiter
}

func Setup() *fiber.App {
	app := fiber.New()
	authKeyGen := func(c *fiber.Ctx) string {
		usr := c.FormValue("username")
		psw := c.FormValue("password")
		encP, _ := auth.HashPassword(psw)
		psw256 := sha256.Sum256([]byte(encP))
		usr256 := sha256.Sum256([]byte(usr))
		key := hex.EncodeToString(psw256[:]) + hex.EncodeToString(usr256[:])
		return key
	}
	defaultKeyGen := func(c *fiber.Ctx) string {
		return utils.CopyString(c.Path())
	}
	authCache := cacheSetupPost(authKeyGen)
	app.Post("/login", authCache, limiterSetup(10), corsSetup("POST"), handlers.HandleLogin)
	app.Post("/register", authCache, limiterSetup(10), corsSetup("POST"), handlers.HandleSignUp)
	defaultCache := cacheSetupGet(defaultKeyGen)
	app.Post("/logout", limiterSetup(10), corsSetup("POST"), handlers.HandleLogout)
	app.Get("/signin", defaultCache, corsSetup("GET"), handlers.LoginRoute)
	app.Get("/signup", defaultCache, corsSetup("GET"), handlers.SignUpRoute)
	app.Get("/categories", corsSetup("GET"), handlers.CategoriesRoute)
	app.Post("/rules", limiterSetup(10), corsSetup("POST"), handlers.HandleCreateRule)
	app.Patch("/rules", limiterSetup(10), corsSetup("POST"), handlers.HandleUpdateRule)
	app.Delete("/rules/:id", limiterSetup(10), corsSetup("DELETE"), handlers.HandleDeleteRule)
	app.Get("/notes", corsSetup("GET"), handlers.FilesRoute)
	app.Post("/notes", limiterSetup(10), corsSetup("POST"), handlers.HandleUploadFile)
	app.Delete("/notes/:id", limiterSetup(10), corsSetup("DELETE"), handlers.HandleDeleteFile)
	app.Get("/", handlers.HomeRoute)
	app.Static("/static", "./static/")
	app.Use(handlers.PageDoesNotExistRoute)
	return app
}

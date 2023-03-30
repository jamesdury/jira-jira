package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/template/html"

	"github.com/jxxxxxxxy/jira-jira/internal/routes"
)

func main() {
	engine := html.New("./web/template", ".html")

	engine.AddFunc("timestamp", func(name string) string {

		// TODO remove last 3 characters from string 000
		s := strings.Split(name, "T")
		return s[0]
	})

	engine.AddFunc("date", func(name string) string {
		s := strings.Split(name, "T")
		return s[0]
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			if c.Path() == "/issues" || c.Path() == "/image" {
				return true
			}

			return c.Query("refresh") == "true"
		},
		Expiration:   30 * time.Minute,
		CacheControl: true,
	}))

	app.Use(logger.New())
	routes.SetupRoutes(app)

	app.Use(recover.New())
	app.Use(favicon.New(favicon.Config{
		File: "./web/static/media/favicon.ico",
	}))

	app.Use(filesystem.New(filesystem.Config{
		Root:   http.Dir("./web"),
		Browse: true,
	}))

    // use a weird port so you can developer at the same time
	log.Fatal(app.Listen(":7070"))
}

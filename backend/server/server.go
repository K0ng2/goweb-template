package server

import (
	"io/fs"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/swagger/v2"

	"project/database"
	"project/docs"
	"project/handler"
	"project/web"
)

// @BasePath /api
func New(db *database.Database) *fiber.App {
	app := fiber.New()

	app.Use(cors.New())

	// Configure logger middleware with settings from Viper config
	loggerConfig := logger.Config{
		Format: logger.CommonFormat,
	}
	app.Use(logger.New(loggerConfig))
	app.Use(recover.New())

	h := handler.NewHandler(db)
	api := app.Group("api")
	api.Get("databasez", h.DatabaseHealth)

	api.Get("/swagger/*", swagger.New(swagger.Config{
		InstanceName:             docs.SwaggerInfo.InstanceName(),
		DefaultModelsExpandDepth: -1,
	}))

	fSys, err := fs.Sub(web.EmbeddedFiles, "public")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
	}

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendFile("index.html", fiber.SendFile{
			FS:       fSys,
			Compress: true,
		})
	})

	app.Get("200", func(c fiber.Ctx) error {
		return c.SendFile("200.html", fiber.SendFile{
			FS:       fSys,
			Compress: true,
		})
	})

	app.Use(static.New("", static.Config{
		FS:       fSys,
		Compress: true,
		NotFoundHandler: func(c fiber.Ctx) error {
			return c.SendFile("404.html", fiber.SendFile{
				FS:       fSys,
				Compress: true,
			})
		},
	}))

	return app
}

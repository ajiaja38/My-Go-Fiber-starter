package main

import (
	"learn/fiber/pkg/err"
	"learn/fiber/pkg/model"
	"learn/fiber/pkg/router"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	logrus.SetLevel(logrus.DebugLevel)

	err := godotenv.Load()

	if err != nil {
		logrus.Error("Error loading .env file")
	}
}

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: err.ErrorHandler,
	})

	port := ":" + os.Getenv("PORT")

	if port == "" {
		port = ":3000"
	}

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(csrf.New())
	app.Use(recover.New())

	route := app.Group("/api/v1")

	route.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(model.ResponseEntity[any]{
			Code:    fiber.StatusOK,
			Message: "Halo Fiber Framework 🚀",
		})
	})

	route.Get("/metrics", monitor.New(monitor.Config{Title: "Fiber Metrics Page"}))

	router.MainRouter(route)

	logrus.Infof("Server running on port http://localhost%s 🚀", port)
	logrus.Fatal(app.Listen(port))
}

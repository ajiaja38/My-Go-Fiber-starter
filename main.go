package main

import (
	"learn/fiber/config"
	_ "learn/fiber/docs"
	"learn/fiber/pkg/err"
	"learn/fiber/pkg/handler"
	"learn/fiber/pkg/repository"
	"learn/fiber/pkg/router"
	"learn/fiber/pkg/service"
	"learn/fiber/utils"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
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

//	@title			         Swagger Fiber API Docs
//	@version		       1.0
//	@description	   Documentation API Fiber By M. Aji Perdana.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	  M. Aji Perdana
//	@contact.email	 ajicooljazz38@gmail.com

//	@license.name	  Apache 2.0
//	@license.url	   http://www.apache.org/licenses/LICENSE-2.0.html

// @host						       localhost:3001
// @BasePath					   /api/v1
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @type						apiKey
// @description				Masukkan token JWT Anda di sini. Contoh: "Bearer <token>"
func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: err.ErrorHandler,
	})

	port := ":" + os.Getenv("PORT")

	if port == "" {
		port = ":3000"
	}

	db := config.DBConfig()

	// Init Repository
	userRepository := repository.NewUserRepository(db)

	// Init Service
	userService := service.NewUserService(userRepository)
	fileService, err := service.NewFileService()

	if err != nil {
		log.Fatalf("Error creating file service: %v", err)
	}

	// Init Handler
	userHandler := handler.NewUserHandler(userService)
	fileHandler := handler.NewFileHandler(fileService)

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(recover.New())

	app.Get("/swagger/*", swagger.HandlerDefault)

	route := app.Group("/api/v1")

	route.Get("/", RootHandler)
	route.Get("/metrics", monitor.New(monitor.Config{Title: "Fiber Metrics Page"}))

	// Init Router
	router.UserRouter(route, userHandler)
	router.FileRouter(route, fileHandler)

	logrus.Infof("Server running on http://localhost%s/api/v1 🚀", port)
	logrus.Fatal(app.Listen(port))
}

// @Summary		    Root Endpoint
// @Description	Returns a welcome message
// @Tags			       status
// @Accept			     json
// @Produce		    json
// @Router			     / [get]
func RootHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse[any](c, fiber.StatusOK, "Halo Fiber Framework 🚀", nil)
}

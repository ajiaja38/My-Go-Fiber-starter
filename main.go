package main

import (
	"learn/fiber/config"
	_ "learn/fiber/docs"
	"learn/fiber/pkg/err"
	"learn/fiber/pkg/handler"
	"learn/fiber/pkg/middleware"
	"learn/fiber/pkg/repository"
	"learn/fiber/pkg/router"
	"learn/fiber/pkg/service"
	"learn/fiber/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

//	@title			         Swagger Fiber API Docs
//	@version		       1.0
//	@description	   Documentation API Fiber By M. Aji Perdana | 2025.
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
	if err := config.LoadEnv(); err != nil {
		log.Errorf("Failed to load environment variables: %v", err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: err.ErrorHandler,
	})

	port := ":" + config.PORT.GetValue()

	if port == "" {
		port = ":3000"
	}

	db := config.DBConfig()

	// Init Repository
	userRepository := repository.NewUserRepository(db)
	blogRepository := repository.NewBlogRepository(db)

	// Init Service
	userService := service.NewUserService(userRepository)
	blogService := service.NewBlogService(blogRepository, userRepository)
	fileService, err := service.NewFileService()

	if err != nil {
		log.Fatalf("Error creating file service: %v", err)
	}

	// Init Handler
	userHandler := handler.NewUserHandler(userService)
	blogHandler := handler.NewBlogHandler(blogService)
	fileHandler := handler.NewFileHandler(fileService)

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(recover.New())
	app.Use(middleware.LimitUploadSize())
	app.Use(limiter.New(limiter.Config{
		Max:        30,
		Expiration: 10 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return fiber.NewError(fiber.StatusTooManyRequests, "Sorry, To Many Request")
		},
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	route := app.Group("/api/v1")

	route.Get("/", RootHandler)
	route.Get("/metrics", monitor.New(monitor.Config{Title: "Fiber Metrics Page"}))

	// Init Router
	router.UserRouter(route, userHandler)
	router.BlogRouter(route, blogHandler)
	router.FileRouter(route, fileHandler)

	log.Infof("Server running on http://127.0.0.1%s/api/v1 🚀", port)
	log.Fatal(app.Listen(port))
}

// @Summary		    Root Endpoint
// @Description	Returns a welcome message
// @Tags			       status
// @Accept			     json
// @Produce		    json
// @Router			     / [get]
func RootHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse[any](c, fiber.StatusOK, "Halo From Fiber 🚀", nil)
}

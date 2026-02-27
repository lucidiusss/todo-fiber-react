package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/lucidiusss/todo-fiber-react/database"
	"github.com/lucidiusss/todo-fiber-react/handlers"
	"github.com/lucidiusss/todo-fiber-react/middleware"
)

func main() {
	// load env variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: .env file not found, using system environment variables")
	}

	// connect to db
	database.ConnectDB()

	// create fiber app
	app := fiber.New(fiber.Config{
		AppName: "Tasks API with Fiber, PostgreSQL and JWT Auth v1",
	})

	// middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// routing
	setupRoutes(app)

	// start server
	app.Listen(":3000")
}

func setupRoutes(app *fiber.App) {
	loginLimiter := limiter.Config{
		Max:        5,
		Expiration: 5 * time.Minute,
		KeyGenerator: func(c fiber.Ctx) string {
			// Rate limit by IP + username (extract username from request)
			var body struct {
				Username string `json:"username"`
			}
			if err := c.Bind().Body(&body); err == nil {
				return "login:" + c.IP() + ":" + body.Username
			}
			return "login:" + c.IP()
		},
		LimitReached: func(c fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success":     false,
				"message":     "Too many login attempts. Please try again later.",
				"retry_after": 300, // 15 minutes in seconds
			})
		},
		SkipFailedRequests: false, // Only count failed attempts
		Storage:            nil,   // Uses memory storage by default
	}
	// public routes
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})
	auth := app.Group("/api/v1/auth")
	auth.Post("/register", limiter.New(loginLimiter), handlers.Register)
	auth.Post("/login",
		limiter.New(loginLimiter),
		handlers.Login,
	)

	// protected routes
	api := app.Group("/api/v1", middleware.AuthRequired())
	tasksRoutes := api.Group("/tasks")
	api.Get("/user", handlers.GetMe)

	tasksRoutes.Get("/", handlers.GetTasks)
	tasksRoutes.Post("/", handlers.CreateTask)
	tasksRoutes.Delete("/:id", handlers.DeleteTask)
	tasksRoutes.Get("/:id", handlers.GetTask)
	tasksRoutes.Put("/:id", handlers.UpdateTask)

}

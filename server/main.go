package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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
	// public routes
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})
	auth := app.Group("/api/v1/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)

	// protected routes
	api := app.Group("/api/v1", middleware.AuthRequired())
	tasksRoutes := api.Group("/tasks")
	api.Get("/profile", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id":    c.Locals("user_id"),
			"created_at": c.Locals("created_at"),
			"username":   c.Locals("username"),
			"message":    "Welcome to your profile!",
		})
	})

	tasksRoutes.Get("/", handlers.GetTasks)
	tasksRoutes.Post("/", handlers.CreateTask)
	tasksRoutes.Delete("/:id", handlers.DeleteTask)
	tasksRoutes.Get("/:id", handlers.GetTask)
	tasksRoutes.Put("/:id", handlers.UpdateTask)

}

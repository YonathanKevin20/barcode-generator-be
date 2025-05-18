package main

import (
	"barcode-generator-be/config"
	"barcode-generator-be/handlers"
	"barcode-generator-be/middleware"
	"barcode-generator-be/migrations"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Redis
	if err := config.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer config.RedisClient.Close()

	// Add a short delay before connecting to the database
	time.Sleep(2 * time.Second)

	// Connect to database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Run migrations
	if err := migrations.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})
	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	// Set up routes
	setupRoutes(app)

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	// Start server
	if err := app.Listen(":8080"); err != nil {
		log.Panic(err)
	}
}

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api", fiber.StatusTemporaryRedirect)
	})

	// Routes
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to the Barcode Generator API",
		})
	})

	// Public routes
	api.Post("/auth/login", handlers.Login)
	// api.Post("/auth/register", handlers.Register)

	// Protected routes
	api.Use(middleware.AuthMiddleware)

	// User routes (admin only)
	users := api.Group("/users")
	users.Use(middleware.AdminOnly)
	users.Get("/", handlers.GetUsers)
	users.Get("/:id", handlers.GetUser)
	users.Post("/", handlers.CreateUser)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)

	// Logout route
	api.Post("/auth/logout", handlers.Logout)

	// Get current user route
	api.Get("/auth/me", handlers.GetMe)

	// Status routes
	api.Get("/statuses", handlers.GetStatuses)

	// Category routes
	api.Get("/categories", handlers.GetCategories)
	api.Get("/categories/:id", handlers.GetCategory)
	api.Post("/categories", handlers.CreateCategory)
	api.Put("/categories/:id", handlers.UpdateCategory)
	api.Delete("/categories/:id", handlers.DeleteCategory)

	// Supplier routes
	api.Get("/suppliers", handlers.GetSuppliers)
	api.Get("/suppliers/:id", handlers.GetSupplier)
	api.Post("/suppliers", handlers.CreateSupplier)
	api.Put("/suppliers/:id", handlers.UpdateSupplier)
	api.Delete("/suppliers/:id", handlers.DeleteSupplier)

	// Barcode routes
	api.Get("/barcodes", handlers.GetBarcodes)
	api.Get("/barcodes/:id", handlers.GetBarcode)
	api.Post("/barcodes", handlers.CreateBarcode)
	api.Delete("/barcodes/:id", handlers.DeleteBarcode)
}

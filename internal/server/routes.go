package server

import (
	"sync"
	"venturo-core/configs"
	"venturo-core/internal/handler/http"
	"venturo-core/internal/middleware"
	"venturo-core/internal/service"
	"venturo-core/internal/adapter/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

func registerRoutes(app *fiber.App, db *gorm.DB, conf *configs.Config, wg *sync.WaitGroup) {
	app.Static("/public", "./public")
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Welcome to Venturo Core!",
		})
	})

	api := app.Group("/api/v1")

	// --- Setups ---
	authMiddleware := middleware.NewAuthMiddleware(conf.JWTSecretKey)

	// --- Setup Adapters ---
	localUploader := storage.NewLocalUploaderAdapter("./public/uploads")

	// --- Setup services ---
	authService := service.NewAuthService(db, conf)
	userService := service.NewUserService(db, wg)
	postService := service.NewPostService(db)

	// Add our new service
	transactionService := service.NewTransactionService(db, wg) // <-- ADD THIS
	// Inject all three dependencies into the ProductService.
	productService := service.NewProductService(db, wg, localUploader)

	// --- Setup handlers ---
	authHandler := http.NewAuthHandler(authService)
	userHandler := http.NewUserHandler(userService)
	postHandler := http.NewPostHandler(postService)
	transactionHandler := http.NewTransactionHandler(transactionService)
	productHandler := http.NewProductHandler(productService)


	// Add our new handler
	// transactionHandler := http.NewTransactionHandler(transactionService) // <-- ADD THIS


	// --- Auth routes ---
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)

	// --- User routes ---
	api.Get("/profile", authMiddleware, userHandler.GetProfile)
	api.Put("/profile", authMiddleware, userHandler.UpdateProfile)

	// --- Register Post Routes ---
	postRoutes := api.Group("/posts")
	postRoutes.Get("/", postHandler.GetAllPosts)                      // Public
	postRoutes.Get("/:id", postHandler.GetPostByID)                   // Public
	postRoutes.Post("/", authMiddleware, postHandler.CreatePost)      // Protected
	postRoutes.Put("/:id", authMiddleware, postHandler.UpdatePost)    // Protected
	postRoutes.Delete("/:id", authMiddleware, postHandler.DeletePost) // Protected

	// --- Transaction routes ---
	api.Post("/transactions", authMiddleware, transactionHandler.CreateTransaction) // <-- ADD THIS
	api.Post("/transactions/:id/pay", authMiddleware, transactionHandler.MarkAsPaid) // <-- ADD THIS
	// --- Product routes ---
	api.Post("/products", authMiddleware, productHandler.CreateProduct)
}

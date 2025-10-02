package server

import (
	"sync"
	"venturo-core/configs"
	"venturo-core/internal/handler/http"
	"venturo-core/internal/middleware"
	"venturo-core/internal/service"
	"venturo-core/internal/adapter/storage"
	"venturo-core/internal/adapter/email"
	"venturo-core/internal/adapter/rabbitmq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func registerRoutes(app *fiber.App, db *gorm.DB, conf *configs.Config, wg *sync.WaitGroup) {
	app.Use(cors.New(cors.Config{
        AllowOrigins: "http://localhost:4000",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
    }))
	
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

	// --- Setup RabbitMQ & Email ---
	rabbitAdapter, err := rabbitmq.NewRabbitMQAdapter("amqp://guest:guest@localhost:5672/", "email_queue")
	if err != nil {
		panic(err)
	}

	emailAdapter := email.NewConsoleEmailAdapter()
	emailService := service.NewEmailService(rabbitAdapter, emailAdapter)

	// Start consumer in background
	wg.Add(1)
	go func() {
		defer wg.Done()
		emailService.StartConsumer()
	}()

	// --- Setups ---
	authMiddleware := middleware.NewAuthMiddleware(conf.JWTSecretKey)

	// --- Setup Adapters ---
	localUploader := storage.NewLocalUploaderAdapter("./public/uploads")

	// --- Setup services ---
	authService := service.NewAuthService(db, conf)
	userService := service.NewUserService(db, wg)
	postService := service.NewPostService(db)

	// User Point Service 
	userPointService := service.NewUserPointService(db)
	// Voucher Service
	voucherService := service.NewVoucherService(db)
	// Add our new service
	transactionService := service.NewTransactionService(db, wg) // <-- ADD THIS
	// Inject all three dependencies into the ProductService.
	productService := service.NewProductService(db, wg, localUploader)
	// Add Inventory Service
	inventoryService := service.NewInventoryService(db, wg)
	// Add Report Service
	reportService := service.NewReportService(db)

	// --- Setup handlers ---
	authHandler := http.NewAuthHandler(authService)
	userHandler := http.NewUserHandler(userService)
	postHandler := http.NewPostHandler(postService)
	transactionHandler := http.NewTransactionHandler(transactionService)
	productHandler := http.NewProductHandler(productService)
	inventoryHandler := http.NewInventoryHandler(inventoryService)
	reportHandler := http.NewReportHandler(reportService)
	userPointHandler := http.NewUserPointHandler(userPointService)
	voucherHadnler := http.NewVoucherHandler(voucherService)
	emailHandler := http.NewEmailHandler(emailService)


	// Add our new handler
	// transactionHandler := http.NewTransactionHandler(transactionService) // <-- ADD THIS


	// --- Auth routes ---
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)
	api.Post("/refresh", authHandler.Refresh)
	

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
	api.Post("/transactions", authMiddleware, transactionHandler.CreateTransaction)
	api.Post("/transactions/:id/pay", authMiddleware, transactionHandler.MarkAsPaid) 
	// --- Product routes ---
	api.Post("/products", authMiddleware, productHandler.CreateProduct)
	api.Get("/products", authMiddleware, productHandler.GetAllProducts)
	// Product routes get by id
	api.Get("/products/:id", authMiddleware, productHandler.GetProductByID)
	// Product routes get by slug
	api.Get("/product/:slug", authMiddleware, productHandler.GetProductBySlug)
	// --- Inventory routes ---
	api.Post("/inventories", authMiddleware ,inventoryHandler.StockIn)
	// --- Report routes---
	api.Get("/inventories/report", authMiddleware, reportHandler.GetInventoryReport)
	// --- User Point route ---
	api.Get("/users/points", authMiddleware, userPointHandler.GetUserPointByUserID)
	api.Post("/users/points", authMiddleware, userPointHandler.CreateUserPoint)
	api.Put("/users/:id/points", authMiddleware, userPointHandler.UpdateUserPoint)
	// --- Voucher routes ---
	api.Post("/vouchers", authMiddleware, voucherHadnler.CreateVoucher)
	api.Post("/vouchers/code", authMiddleware, voucherHadnler.GetVoucherByCode)
	// --- Email routes ---
	api.Post("/emails", emailHandler.SendEmail)

}

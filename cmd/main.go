package main

import (
	"fmt"
	"log"
	"os"

	"E-book-service/internal/domain"
	"E-book-service/internal/handler"
	"E-book-service/internal/middleware"
	"E-book-service/internal/repository"
	"E-book-service/internal/service"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title E-book Service API
// @version 1.0
// @description API сервер для управления электронными книгами.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_SSLMODE"),
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Автомиграция
	if err := db.AutoMigrate(&domain.User{}, &domain.Author{}, &domain.Book{}, &domain.Review{}, &domain.Shelf{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	repo := repository.NewRepository(db)
	svc := service.NewService(repo, jwtSecret)
	h := handler.NewHandler(svc)

	e := echo.New()
	e.Use(echoMW.Recover())
	e.Use(middleware.RateLimiter(rdb))
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Routes (PUBLIC)
	e.GET("/health", h.Health)
	e.POST("/register", h.Register)
	e.POST("/login", h.Login)

	// Routes (PROTECTED)

	a := e.Group("/api/v1")
	a.Use(middleware.JWTMiddleware(jwtSecret))

	{
		a.GET("/me", h.GetMe)
		a.PUT("/me", h.UpdateProfile)
		a.GET("/profile", h.GetMe)

		// Books
		a.GET("/books", h.ListBooks)
		a.POST("/books", h.CreateBook)
		a.GET("/books/:id", h.GetBook)
		a.PUT("/books/:id", h.UpdateBook)
		a.DELETE("/books/:id", h.DeleteBook)
		a.GET("/books/:id/content", h.GetBookContent)

		// Authors
		a.GET("/authors", h.ListAuthors)
		a.POST("/authors", h.CreateAuthor)
		a.GET("/authors/:id", h.GetAuthor)
		a.PUT("/authors/:id", h.UpdateAuthor)
		a.DELETE("/authors/:id", h.DeleteAuthor)
		a.GET("/authors/:id/books", h.GetAuthorBooks)

		// Reviews
		a.GET("/books/:id/reviews", h.ListReviews)
		a.POST("/books/:id/reviews", h.AddReview)
		a.DELETE("/reviews/:id", h.DeleteReview)

		// Shelf
		a.GET("/shelf", h.GetShelf)
		a.POST("/shelf/:id", h.AddToShelf)
		a.DELETE("/shelf/:id", h.RemoveFromShelf)
		a.PUT("/shelf/:id", h.AddToShelf)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	log.Printf("Server starting on %s", port)
	e.Logger.Fatal(e.Start(port))
}

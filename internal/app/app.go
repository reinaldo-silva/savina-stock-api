package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/reinaldo-silva/savina-stock/config"
	"github.com/reinaldo-silva/savina-stock/internal/domain/auth"
	"github.com/reinaldo-silva/savina-stock/internal/domain/category"
	"github.com/reinaldo-silva/savina-stock/internal/domain/image_service"
	"github.com/reinaldo-silva/savina-stock/internal/domain/product"
	"github.com/reinaldo-silva/savina-stock/internal/domain/product_image"
	"github.com/reinaldo-silva/savina-stock/internal/domain/user"
	"github.com/reinaldo-silva/savina-stock/internal/infrastructure/db/gorm"
	s3_provider "github.com/reinaldo-silva/savina-stock/internal/infrastructure/image_provider/aws"
	jwt_middleware "github.com/reinaldo-silva/savina-stock/internal/middleware/jwt"
)

type App struct {
	Router *chi.Mux
}

func (a *App) Initialize(cfg *config.Config) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	connection := gorm.NewGormDB(dsn)

	s3Config := config.LoadS3Config()

	s3Provider, err := s3_provider.NewS3Provider(s3Config)
	if err != nil {
		log.Fatal("failed to initialize cloudinary service: ", err)
	}

	a.Router = chi.NewRouter()

	isProduction := os.Getenv("ENVIRONMENT") == "production"

	var allowedOrigins []string
	if isProduction {
		allowedOrigins = []string{fmt.Sprintf("https://%s", os.Getenv("HOST_WEB")), fmt.Sprintf("https://www.%s", os.Getenv("HOST_WEB"))}
	} else {
		allowedOrigins = []string{"http://localhost:3000"}
	}

	jwtMiddleware := jwt_middleware.NewJwtMiddleware([]byte(cfg.JwtSecret))

	a.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)

	userRepo := gorm.NewGormUserRepository(connection)
	productRepo := gorm.NewGormProductRepository(connection)
	categoryRepo := gorm.NewCategoryRepository(connection)
	imageRepo := gorm.NewGormImageRepository(connection)

	imageService := image_service.NewImageService(s3Provider)

	userUseCase := user.NewUserUseCase(userRepo)
	productUseCase := product.NewProductUseCase(productRepo, categoryRepo, imageRepo, imageService)
	categoryUseCase := category.NewCategoryUseCase(categoryRepo)
	imageUseCase := product_image.NewImageUseCase(imageService, imageRepo)

	authHandler := auth.NewAuthHandler(userUseCase)
	userHandler := user.NewUserHandler(userUseCase)
	productHandler := product.NewProductHandler(productUseCase, imageService)
	categoryHandler := category.NewCategoryHandler(categoryUseCase)
	imageHandler := product_image.NewImageHandler(imageUseCase)

	a.Router.Route("/users", func(r chi.Router) {
		r.Use(jwtMiddleware.ValidateToken)
		r.Use(jwtMiddleware.RequireRoles(string(user.AdminRole)))
		r.Get("/", userHandler.GetUsers)
		r.Post("/", userHandler.CreateUser)
	})

	a.Router.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", authHandler.SignUp)
		r.Post("/sign-in", authHandler.SignIn)
	})

	a.Router.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.GetProducts)
		r.Get("/{slug}", productHandler.GetProductBySlug)
		r.Get("/{slug}/images", productHandler.GetProductImages)
		r.Group(func(r chi.Router) {
			r.Use(jwtMiddleware.ValidateToken)
			r.Use(jwtMiddleware.RequireRoles(string(user.AdminRole)))
			r.Get("/to-admin", productHandler.GetProductsToAdmin)
			r.Get("/to-admin/{slug}", productHandler.GetProductBySlugToAdmin)
			r.Post("/", productHandler.CreateProduct)
			r.Delete("/{slug}", productHandler.DeleteProduct)
			r.Put("/{slug}", productHandler.UpdateProduct)
			r.Patch("/{slug}/upload-image", productHandler.UploadImages)
			r.Patch("/{slug}/categories/link", productHandler.LinkCategories)
			r.Patch("/{slug}/cover/{uuid}", imageHandler.SetImageAsCover)
		})

	})

	a.Router.Route("/category", func(r chi.Router) {
		r.Get("/", categoryHandler.GetAllCategories)
		r.Get("/{id}", categoryHandler.GetCategoryByID)
		r.Group(func(r chi.Router) {
			r.Use(jwtMiddleware.ValidateToken)
			r.Use(jwtMiddleware.RequireRoles(string(user.AdminRole)))
			r.Post("/", categoryHandler.CreateCategory)
			r.Delete("/{id}", categoryHandler.DeleteCategory)
			r.Put("/{id}", categoryHandler.UpdateCategory)
		})
	})

	a.Router.Route("/image", func(r chi.Router) {
		r.Get("/{uuid}", imageHandler.GetImage)
		r.Group(func(r chi.Router) {
			r.Use(jwtMiddleware.ValidateToken)
			r.Use(jwtMiddleware.RequireRoles(string(user.AdminRole)))
			r.Delete("/{uuid}", imageHandler.DeleteImage)
		})
	})

	a.Router.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

}

func (a *App) Run(cfg *config.Config) {
	fmt.Printf("Server running on port %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, a.Router))
}

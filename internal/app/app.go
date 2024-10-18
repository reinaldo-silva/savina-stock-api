package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	api_category "github.com/reinaldo-silva/savina-stock/api/category"
	api_image "github.com/reinaldo-silva/savina-stock/api/image"
	api_product "github.com/reinaldo-silva/savina-stock/api/product"
	"github.com/reinaldo-silva/savina-stock/config"
	"github.com/reinaldo-silva/savina-stock/internal/domain/category"
	"github.com/reinaldo-silva/savina-stock/internal/domain/image"
	"github.com/reinaldo-silva/savina-stock/internal/domain/product"
	s3_provider "github.com/reinaldo-silva/savina-stock/internal/provider/aws"
	category_repository "github.com/reinaldo-silva/savina-stock/internal/repository/category"
	image_repository "github.com/reinaldo-silva/savina-stock/internal/repository/image"
	product_repository "github.com/reinaldo-silva/savina-stock/internal/repository/product"
	service "github.com/reinaldo-silva/savina-stock/internal/service/image"
	usecase_category "github.com/reinaldo-silva/savina-stock/internal/usecase/category"
	usecase_image "github.com/reinaldo-silva/savina-stock/internal/usecase/image"
	usecase_product "github.com/reinaldo-silva/savina-stock/internal/usecase/product"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Router *chi.Mux
	DB     *gorm.DB
}

func (a *App) Initialize(cfg *config.Config) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	var err error
	a.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	err = a.DB.AutoMigrate(
		&product.Product{},
		&image.ProductImage{},
		&category.Category{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	// cloudinaryConfig := config.LoadCloudinaryConfig()
	s3Config := config.LoadS3Config()

	// cloudinaryProvider, err := provider.NewCloudinaryProvider(cloudinaryConfig)
	// if err != nil {
	// 	log.Fatal("failed to initialize cloudinary service: ", err)
	// }

	s3Provider, err := s3_provider.NewS3Provider(s3Config)
	if err != nil {
		log.Fatal("failed to initialize cloudinary service: ", err)
	}

	a.Router = chi.NewRouter()

	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.Recoverer)

	productRepo := product_repository.NewGormProductRepository(a.DB)
	categoryRepo := category_repository.NewCategoryRepository(a.DB)
	imageRepo := image_repository.NewGormImageRepository(a.DB)

	imageService := service.NewImageService(s3Provider)

	productUseCase := usecase_product.NewProductUseCase(productRepo, imageRepo)
	categoryUseCase := usecase_category.NewCategoryUseCase(categoryRepo)
	imageUseCase := usecase_image.NewImageUseCase(imageService)

	productHandler := api_product.NewProductHandler(productUseCase, imageService)
	categoryHandler := api_category.NewCategoryHandler(categoryUseCase)
	imageHandler := api_image.NewImageHandler(imageUseCase)

	a.Router.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.GetProducts)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{slug}", productHandler.GetProductBySlug)
		r.Delete("/{slug}", productHandler.DeleteProduct)
		r.Put("/{slug}", productHandler.UpdateProduct)
		r.Patch("/{slug}/upload-image", productHandler.UploadImages)
		r.Get("/{slug}/images", productHandler.GetProductImages)
	})

	a.Router.Route("/category", func(r chi.Router) {
		r.Post("/", categoryHandler.CreateCategory)
		r.Get("/", categoryHandler.GetAllCategories)
	})

	a.Router.Route("/image", func(r chi.Router) {
		r.Get("/{uuid}", imageHandler.GetImage)
	})

}

func (a *App) Run(cfg *config.Config) {
	fmt.Printf("Server running on port %s\n", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, a.Router))
}

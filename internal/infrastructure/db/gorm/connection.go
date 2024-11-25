package gorm

import (
	"log"

	"github.com/reinaldo-silva/savina-stock/internal/domain/category"
	"github.com/reinaldo-silva/savina-stock/internal/domain/product"
	"github.com/reinaldo-silva/savina-stock/internal/domain/product_audit"
	"github.com/reinaldo-silva/savina-stock/internal/domain/product_image"
	"github.com/reinaldo-silva/savina-stock/internal/domain/sale"
	"github.com/reinaldo-silva/savina-stock/internal/domain/sale_item"
	"github.com/reinaldo-silva/savina-stock/internal/domain/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(dsn string) *gorm.DB {
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	err = connection.AutoMigrate(
		&product.Product{},
		&product_image.ProductImage{},
		&category.Category{},
		&user.User{},
		&product_audit.ProductAudit{},
		&sale.Sale{},
		&sale_item.SaleItem{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	return connection

}

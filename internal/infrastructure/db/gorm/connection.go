package gorm

import (
	"log"

	"github.com/reinaldo-silva/savina-stock/internal/domain/category"
	"github.com/reinaldo-silva/savina-stock/internal/domain/product"
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
		// &domain.ProductImage{},
		&category.Category{},
		&user.User{})
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}

	return connection

}

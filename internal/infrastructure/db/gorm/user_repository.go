package gorm

import (
	"github.com/reinaldo-silva/savina-stock/internal/domain/user"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) user.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) GetAll(
	page int,
	pageSize int) ([]user.User, int64, error) {
	var users []user.User
	var total int64

	query := r.db

	if err := query.Model(&user.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *GormUserRepository) Create(p user.User) error {
	return r.db.Create(&p).Error
}

func (r *GormUserRepository) FindByEmail(email string) (*user.User, error) {
	var user user.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

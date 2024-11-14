package user

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	AdminRole  Role = "ADMIN"
	ClientRole Role = "CLIENT"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(150);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password"`
	Role      Role      `gorm:"type:varchar(20);not null;default:CLIENT" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserRepository interface {
	GetAll(page int,
		pageSize int) ([]User, int64, error)
	Create(user User) error
	FindByEmail(email string) (*User, error)
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	return u.validateRole()
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	return u.validateRole()
}

func (u *User) validateRole() error {
	if u.Role == "" {
		return nil
	}

	if u.Role != AdminRole && u.Role != ClientRole {
		return fmt.Errorf("invalid role: %s", u.Role)
	}
	return nil
}

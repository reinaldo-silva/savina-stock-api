package usecase_user

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/reinaldo-silva/savina-stock/config"
	"github.com/reinaldo-silva/savina-stock/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo user.UserRepository
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewUserUseCase(repo user.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (uc *UserUseCase) GetAll(page int, pageSize int) ([]user.UserResponse, int64, error) {
	users, total, err := uc.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	var userResponses []user.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, *user.ToResponse())
	}

	return userResponses, total, nil
}

func (uc *UserUseCase) Create(u user.User) (*user.User, error) {

	if strings.TrimSpace(u.Name) == "" {
		return nil, errors.New("user name cannot be empty")
	}

	err := uc.repo.Create(u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (uc *UserUseCase) GetByEmail(email string) (*user.UserResponse, error) {
	user, err := uc.repo.FindByEmail(email)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user.ToResponse(), err
}

func (uc *UserUseCase) SignInUseCase(email string, pass string) (*user.UserResponse, string, error) {

	existingUser, err := uc.repo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(pass))
	if err != nil {
		return nil, "", errors.New("invalid email or password")

	}

	token, err := uc.generateJWT(existingUser)
	if err != nil {
		return nil, "", errors.New("error generating token")
	}

	return existingUser.ToResponse(), token, nil
}

func (uc *UserUseCase) generateJWT(user *user.User) (string, error) {
	cfg := config.LoadConfig()

	claims := Claims{
		UserID: user.ID,
		Role:   string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JwtSecret))
}

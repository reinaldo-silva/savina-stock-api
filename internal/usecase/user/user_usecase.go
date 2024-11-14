package usecase_user

import (
	"errors"
	"strings"

	"github.com/reinaldo-silva/savina-stock/internal/domain/user"
)

type UserUseCase struct {
	repo user.UserRepository
}

func NewUserUseCase(repo user.UserRepository) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (uc *UserUseCase) GetAll(page int, pageSize int) ([]user.User, int64, error) {
	users, total, err := uc.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
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

func (uc *UserUseCase) GetByEmail(email string) (*user.User, error) {
	return uc.repo.FindByEmail(email)
}

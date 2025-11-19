package repo

import (
	"context"
	"errors"

	"github.com/visualect/tl/internal/dto"
	"github.com/visualect/tl/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsersRepository interface {
	Create(ctx context.Context, data dto.RegisterUserRequest) error
}

type usersRepo struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) UsersRepository {
	return &usersRepo{db}
}

func (u *usersRepo) Create(ctx context.Context, data dto.RegisterUserRequest) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		return err
	}

	err = gorm.G[models.User](u.db).Create(ctx, &models.User{Login: data.Login, PasswordHash: string(hashed)})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errors.New("user with this login name already exists")
	}
	return err
}

package repo

import (
	"context"
	"errors"

	"github.com/visualect/tl/internal/dto"
	"github.com/visualect/tl/internal/models"
	"gorm.io/gorm"
)

type UsersRepository interface {
	Create(ctx context.Context, data dto.RegisterUserRequest) error
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}

type usersRepo struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) UsersRepository {
	return &usersRepo{db}
}

func (u *usersRepo) Create(ctx context.Context, data dto.RegisterUserRequest) error {
	err := gorm.G[models.User](u.db).Create(ctx, &models.User{Login: data.Login, PasswordHash: data.Password})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errors.New("User with this login name already exists")
	}
	return err
}

func (u *usersRepo) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	user, err := gorm.G[models.User](u.db).Where("login = ?", login).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("User is not registered")
		}
		return models.User{}, err
	}
	return user, nil
}

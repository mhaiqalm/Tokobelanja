package usecase

import (
	"context"
	"time"

	"Tokobelanja/app/helper"
	"Tokobelanja/domain"
)

type UserUsecase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(u domain.UserRepository) domain.UserUsecase {
	return &UserUsecase{userRepository: u}
}

func (u *UserUsecase) Login(ctx context.Context, user *domain.User) (token string, err error) {
	res, err := u.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		err = domain.ErrUnauthorized
		return
	}
	comparePass := helper.ComparePass([]byte(res.Password), []byte(user.Password))
	if !comparePass {
		err = domain.ErrUnauthorized
		return
	}
	token = helper.GenerateToken(res.ID, res.Role)
	return
}

func (u *UserUsecase) Register(ctx context.Context, user *domain.User) (domain.User, error) {
	_, err := u.userRepository.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return domain.User{}, domain.ErrConflict
	}

	user.Password = helper.HassPass(user.Password)
	user.Role = "customer"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	userId, err := u.userRepository.StoreUser(ctx, user)
	if err != nil {
		return domain.User{}, domain.ErrInternalServerError
	}
	userData, err := u.userRepository.GetUserByID(ctx, userId)
	if err != nil {
		return domain.User{}, domain.ErrNotFound
	}
	return userData, nil
}

func (u *UserUsecase) TopUp(ctx context.Context, user *domain.User) (domain.User, error) {
	userTemp, err := u.userRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		return domain.User{}, domain.ErrNotFound
	}

	user.Balance += userTemp.Balance

	user.UpdatedAt = time.Now()
	err = u.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return domain.User{}, domain.ErrInternalServerError
	}
	return *user, nil
}

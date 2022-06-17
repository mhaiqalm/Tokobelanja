package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement:true"`
	FullName  string    `json:"full_name" gorm:"notNull"`
	Email     string    `json:"email" gorm:"notNull"`
	Password  string    `json:"password" gorm:"notNull"`
	Role      string    `json:"role"`
	Balance   int64     `json:"balance" gorm:"notNull"`
	CreatedAt time.Time `json:"created_at" gorm:"notNull"`
	UpdatedAt time.Time `json:"updated_at" gorm:"notNull"`
}

type UserUsecase interface {
	Login(ctx context.Context, user *User) (token string, err error)
	Register(ctx context.Context, user *User) (User, error)
	TopUp(ctx context.Context, user *User) (User, error)
}

type UserRepository interface {
	StoreUser(ctx context.Context, user *User) (userId int64, err error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	UpdateUser(ctx context.Context, user *User) error
}

package model

import (
	"errors"
	"gorm.io/gorm"
)

type AuthUser struct {
	Email              string `gorm:"primary_key"`
	Username           string
	Password           string
	IsEmailActivated   bool
	RegisterTimestamp  int64
	LastLoginTimestamp int64
	RegisterIp         string
	LastLoginIp        string
	Role               string
}

func (u AuthUser) GetUser(Db *gorm.DB, email string) (user AuthUser, err error) {
	err = Db.Where("email=?", email).Find(&user).Error
	return
}

// CreateUser will register a new user if its email not exist.
func (u *AuthUser) CreateUser(Db *gorm.DB, user *AuthUser) (err error) {
	if _, err = u.GetUser(Db, user.Email); err != nil { //Already registered
		err = errors.New("this email is already registered")
	} else {
		err = Db.Create(user).Error
	}
	return
}

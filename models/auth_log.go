package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Auth_log struct {
	ID        int `gorm:"primary_key" json:"id"`
	CreatedOn int `json:"created_on"`

	Username string `json:"UserName"`
	Userip   string `json:"UserIp`
}

func (auth_log *Auth_log) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())

	return nil
}
func AddAuth_log(username string, userip string) bool {

	fmt.Println("username", username, "userip", userip)
	db.Create(&Auth_log{
		Username: username,
		Userip:   userip,
	})

	return true
}

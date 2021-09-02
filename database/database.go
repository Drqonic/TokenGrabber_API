package database

import (
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DBConn *gorm.DB
)

type User struct {
	gorm.Model
	Hostname string `json:"hostname"`
	UID      string `json:"UID"`
	Token    string `json:"token"`
	Email    string `json:"email"`
}

func Token(user *User, id, mail, token, host string) {
	db := DBConn
	db.Where(User{UID: id}).FirstOrCreate(&user)
	db.First(&user)
	if host != "" {
		user.Hostname = host
	}
	user.Email = mail
	user.Token = token
	db.Save(&user)
}

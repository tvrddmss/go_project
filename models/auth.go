package models

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Img      string `json:"img"`
}

func ExistAuthByID(id int) bool {
	var auth Auth
	db.Select("id").Where("id = ?", id).First(&auth)
	if auth.ID > 0 {
		return true
	}

	return false
}

func CheckAuth(username, password string) bool {
	var auth Auth
	db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth)
	if auth.ID > 0 {
		return true
	}

	return false
}
func Login(username, password string) int {
	var auth Auth
	db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth)
	return auth.ID

}
func GetUserInfo(userid int) (userinfo Auth) {
	db.Where("id = ?", userid).First(&userinfo)
	return

}

func EditAuth(id int, data interface{}) bool {
	db.Model(&Auth{}).Where("id = ?", id).Updates(data)

	return true
}

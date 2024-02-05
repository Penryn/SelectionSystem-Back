package userService

import (
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/config/database"
	"math/rand"
)

func CreateUser(user models.User) error {
	result_a := database.DB.Create(&user)
	result_b := database.DB.Create(&models.Student{UserID: user.ID,StudentID: user.Username,})
	if result_a.Error != nil  {
		return result_a.Error
	}else if result_b.Error != nil {
		return result_b.Error
	}else{
		return nil
	}
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	result := database.DB.Where(models.User{Username: username}).First(&user)
	return user, result.Error
}


func GetAvatar() string {
	avatars := []string{
		"http://inews.gtimg.com/newsapp_bt/0/14710913833/1000",
		"https://i01piccdn.sogoucdn.com/5f97fd70f583cec5",
		"https://c-ssl.duitang.com/uploads/item/201907/23/20190723224932_cykee.thumb.700_0.jpg",
		"http://inews.gtimg.com/newsapp_bt/0/13044084095/1000",
		"https://i02piccdn.sogoucdn.com/9831e5c44cbb5ecc",
		"http://inews.gtimg.com/newsapp_bt/0/13814264984/1000",
	}

	randomIndex := rand.Intn(len(avatars))

	return avatars[randomIndex]
}
package userService

import (
	"github.com/xuri/excelize/v2"

	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/config/config"
	"SelectionSystem-Back/config/database"
	"fmt"
	"math/rand"
)

func CreateUser(user models.User) error {
	result_a := database.DB.Create(&user)
	result_b := database.DB.Create(&models.Student{UserID: user.ID, StudentID: user.Username})
	if result_a.Error != nil {
		return result_a.Error
	} else if result_b.Error != nil {
		return result_b.Error
	} else {
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

// 新建管理员
func CreateAdministrator() error {
	uname := config.Config.GetString("Administrator.Name")
	upass := config.Config.GetString("Administrator.Pass")
	_, err := GetUserByUsername(uname)
	if err == nil {
		return nil
	}
	result := database.DB.Create(&models.User{Username: uname, Password: upass, Type: 3})
	return result.Error
}

// 导入教师excel表
func ImportTeacherExcel() error {
	file, err := excelize.OpenFile("德育导师名单.xlsx")
	if err != nil {
		fmt.Println(1)
		return err
	}
	records, err := file.GetRows("Sheet1")
    if err != nil {
        fmt.Println(2)
        return err
    }
	for i, record := range records {
		if i == 0 || i == 1 {
			continue
		}
		result := database.DB.Create(&models.User{Username: "114514" + record[0], Password: "114514", Type: 2})
		if result.Error != nil {
			return result.Error
		}
		user, err := GetUserByUsername("114514" + record[0])
		if err != nil {
			return err
		}
		fmt.Println(4)
		fmt.Println(record[1], record[2], record[3], record[4], record[5])
		result = database.DB.Create(&models.Teacher{
			UserID:      user.ID,
			TeacherName: record[1],
			Section:     record[2],
			Office:      record[3],
			Phone:       record[4],
			Email:       record[5],
		})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

package userService

import (
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/config/config"
	"SelectionSystem-Back/config/database"
	"math/rand"
)

func CreateUser(user models.User) error {
	result_a := database.DB.Create(&user)
	result_b := database.DB.Omit("teacher_id").Create(&models.Student{UserID: user.ID, StudentID: user.Username})
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

func GetAvartar() string {
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


func GetUserByID(id int) (models.User, error) {
	var user models.User
	result := database.DB.Where(models.User{ID: id}).First(&user)
	return user, result.Error
}

func UpdatePassword(user models.User, newPassword string) error {
	result := database.DB.Model(&user).Update("password", newPassword)
	return result.Error
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
		return err
	}
	records, err := file.GetRows("Sheet1")
    if err != nil {
        return err
    }
	for i, record := range records {
		if i == 0 || i == 1 {
			continue
		}
		result := database.DB.Where(models.User{Username: "114514" + record[0]}).First(&models.User{})
		if result.Error == gorm.ErrRecordNotFound{
			result := database.DB.Create(&models.User{Username: "114514" + record[0], Password: "114514", Type: 2, Avartar: GetAvartar(),})
			if result.Error != nil {
				return result.Error
			}
			user, err := GetUserByUsername("114514" + record[0])
			if err != nil {
				return err
			}
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
		}else if result.Error == nil {
			return nil
		}else{
			return result.Error
		}
	}
	return nil
}

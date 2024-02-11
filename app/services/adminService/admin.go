package adminService

import (
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/config/database"
	"time"

	"gorm.io/gorm"
)

func SetDDL(time time.Time,ddltype,id  int) error {
	var result *gorm.DB
	if ddltype==1{
		result=database.DB.Model(&models.DDL{}).Where(models.DDL{UserID: id,DDLType: 2}).Update("first_time",time)
		return result.Error
	}else if ddltype==2{
		result=database.DB.Model(&models.DDL{}).Where(models.DDL{UserID: id,DDLType: 2}).Update("second_time",time)
		return result.Error
	}
	return result.Error
}


func GetAdvices() ([]models.Advice, error) {
	var advices []models.Advice
	result := database.DB.Find(&advices)
	return advices, result.Error
}


func ResetPassword(user_id int, password string) error {
	result := database.DB.Model(&models.User{ID: user_id}).Update("password", password)
	return result.Error
}
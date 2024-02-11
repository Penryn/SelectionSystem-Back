package teacherService

import (
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/config/database"
)

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	result := database.DB.Where(models.User{
		ID: id,
	}).First(&user)
	return &user, result.Error
}

func StudentList() ([]models.Student, error) {
	var studentList []models.Student
	result := database.DB.Model(models.Student{}).Where(models.Student{
		TargetStatus: 0,
		AdminStatus:  0,
	}).Find(&studentList)
	return studentList, result.Error
}

func StudentCheckList() ([]models.Student, error) {
	var studentList []models.Student
	result := database.DB.Model(models.Student{}).Where("target_status IN (?)", []int{3, 4}).Find(&studentList)
	return studentList, result.Error
}

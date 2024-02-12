package teacherService

import (
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/config/database"
	"gorm.io/gorm"
)

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	result := database.DB.Where(models.User{
		ID: id,
	}).First(&user)
	return &user, result.Error
}

func StudentList(targetId int) ([]models.Student, error) {
	var studentList []models.Student
	result := database.DB.Model(models.Student{}).Where(models.Student{
		TargetID:     targetId,
		TargetStatus: 0,
		AdminStatus:  0,
	}).Find(&studentList)
	return studentList, result.Error
}

func StudentCheckList(targetId int) ([]models.Student, error) {
	var studentList []models.Student
	result := database.DB.Model(models.Student{}).Where(models.Student{
		TargetID: targetId,
	}).Where("target_status IN (?)", []int{1, 2}).Find(&studentList)
	return studentList, result.Error
}

func GetStudentInfoByID(Id int) (*models.Student, error) {
	var info models.Student
	result := database.DB.Where(&models.Student{
		ID: Id,
	}).First(&info)
	if result.Error == gorm.ErrRecordNotFound {
		info.ID = Id
		return &info, result.Error
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &info, nil
}

func UpdateStudentInfo(Id int, info *models.Student) error {
	result := database.DB.Model(models.Student{}).Where(&models.Student{
		ID: Id,
	}).Updates(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

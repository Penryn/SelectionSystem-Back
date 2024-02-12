package studentService

import (
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/config/database"
	"gorm.io/gorm"
)

func StudentExistByPhone(userId int, phone string) error {
	var StudentInfo models.Student
	result := database.DB.Not("user_id = ?", userId).Where(models.Student{
		Phone: phone,
	}).First(&StudentInfo)
	return result.Error
}

func StudentExistByEmail(userId int, email string) error {
	var StudentInfo models.Student
	result := database.DB.Not("user_id = ?", userId).Where(models.Student{
		Email: email,
	}).First(&StudentInfo)
	return result.Error
}

func CreateStudentInfo(userId int, info models.Student) error {
	result := database.DB.Model(models.Student{}).Where(models.Student{
		UserID: userId,
	}).Updates(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAdminDDL() (*models.DDL, error) {
	var adminDDL *models.DDL
	result := database.DB.Model(models.DDL{
		DDLType: 2,
	}).First(&adminDDL)
	if result.Error != nil {
		return nil, result.Error
	}
	return adminDDL, nil
}

func GetStudentInfoByUserID(userId int) (*models.Student, error) {
	var info models.Student
	result := database.DB.Where(&models.Student{
		UserID: userId,
	}).First(&info)
	if result.Error == gorm.ErrRecordNotFound {
		info.UserID = userId
		return &info, result.Error
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &info, nil
}

func UpdateStudentInfo(userId int, info models.Student) error {
	result := database.DB.Model(models.Student{}).Where(&models.Student{
		UserID: userId,
	}).Updates(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateTargetTeacher(userId, targetId int, info *models.Student) error {
	info.TargetID = targetId
	result := database.DB.Model(models.Student{}).Where(&models.Student{
		UserID: userId,
	}).Updates(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetTeacherList(pageNum, pageSize int) ([]models.Teacher, error) {
	var teacherList []models.Teacher
	result := database.DB.Model(models.Teacher{}).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&teacherList)
	return teacherList, result.Error
}

func GetTotalPageNum() (*int64, error) {
	var pageNum int64
	result := database.DB.Model(models.Teacher{}).Count(&pageNum)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pageNum, nil
}

func GetTeacherByTeacherID(teacherId int) (*models.Teacher, error) {
	var teacher *models.Teacher
	result := database.DB.Where(models.Teacher{
		ID: teacherId,
	}).First(&teacher)
	return teacher, result.Error
}

func GetTeacherDDLByUserID(userId int) (models.DDL, error) {
	var ddl models.DDL
	result := database.DB.Where(models.DDL{
		UserID:  userId,
		DDLType: 1,
	})
	return ddl, result.Error
}

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	result := database.DB.Where(models.User{
		ID: id,
	}).First(&user)
	return &user, result.Error
}

func CreateAdvice(advice models.Advice) error {
	result := database.DB.Create(&advice)
	return result.Error
}

func GetAdvice(userId int) ([]models.Advice, error) {
	var advice []models.Advice
	result := database.DB.Model(models.Advice{}).Where(models.Advice{
		UserID: userId,
	}).Find(&advice)
	return advice, result.Error
}

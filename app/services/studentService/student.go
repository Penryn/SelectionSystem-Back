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

func CreateDDLRecord(ddlRecord models.DDL) error {
	result := database.DB.Model(models.DDL{}).Create(&ddlRecord)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAdminDDL(userId int) (*models.DDL, error) {
	var adminDDL *models.DDL
	result := database.DB.Model(models.DDL{
		UserID:  userId,
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

func UpdateStudentInfo(userID int, info models.Student) error {
	result := database.DB.Model(models.Student{}).Where(&models.Student{
		UserID: userID,
	}).Updates(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetTeacherList() ([]models.Teacher, error) {
	var teacherList []models.Teacher
	result := database.DB.Find(&teacherList)
	return teacherList, result.Error
}

func CheckTeacherList(teacher models.Teacher) (int64, error) {
	var studentCount int64
	result := database.DB.Model(&models.Student{}).Where(models.Student{
		TeacherID: teacher.ID,
	}).Count(&studentCount)
	return studentCount, result.Error
}

func PageTeacherList(responseTeacherList []models.Teacher, page_num, page_size int) ([]models.Teacher, int) {
	totalItems := len(responseTeacherList)
	totalPageNum := totalItems / page_size
	if totalItems%page_size != 0 {
		totalPageNum++
	}
	startIndex := (page_num - 1) * page_size
	endIndex := startIndex + page_size
	if endIndex > totalItems {
		endIndex = totalItems
	}
	if startIndex > totalItems {
		startIndex = totalItems
	}
	paginatedList := responseTeacherList[startIndex:endIndex]
	return paginatedList, totalPageNum
}

func GetTeacherByTeacherID(teacherId int) (*models.Teacher, error) {
	var teacher *models.Teacher
	result := database.DB.Where(models.Teacher{
		ID: teacherId,
	}).First(&teacher)
	return teacher, result.Error
}

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	result := database.DB.Where(models.User{
		ID: id,
	}).First(&user)
	return &user, result.Error
}

func UpdateAvatar(userId int, avatar string) error {
	var user *models.User
	if err := database.DB.Where(models.User{
		ID: userId,
	}).First(&user).Error; err != nil {
		return err
	}
	user.Avartar = avatar
	if err := database.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

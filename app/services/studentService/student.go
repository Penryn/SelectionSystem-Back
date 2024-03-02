package studentService

import (
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/utils"
	"SelectionSystem-Back/config/database"
	"gorm.io/gorm"
)

func StudentExistByPhone(userId int, phone string) error {
	var StudentInfo models.Student
	phone = utils.AesEncrypt(phone)
	result := database.DB.Not("user_id = ?", userId).Where(models.Student{
		Phone: phone,
	}).First(&StudentInfo)
	return result.Error
}

func StudentExistByEmail(userId int, email string) error {
	var StudentInfo models.Student
	email = utils.AesEncrypt(email)
	result := database.DB.Not("user_id = ?", userId).Where(models.Student{
		Email: email,
	}).First(&StudentInfo)
	return result.Error
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
	aseDecryptStudentInfo(&info)
	return &info, nil
}

func UpdateStudentInfo(userId int, info models.Student) error {
	aseEncryptStudentInfo(&info)
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
	info.TargetStatus = 1
	info.AdminStatus = 0
	aseEncryptStudentInfo(info)
	result := database.DB.Model(models.Student{}).Where(&models.Student{
		UserID: userId,
	}).Updates(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateTeacher(targetId int, studentsNum int) error {
	result := database.DB.Model(&models.Teacher{}).Where("id = ?", targetId).UpdateColumn("students_num", studentsNum)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetTeacherList(pageNum, pageSize int,name string) ([]models.Teacher, error) {
	var teacherList []models.Teacher
	result := database.DB.Where(models.Teacher{
		TeacherName: name,
	}).Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&teacherList)
	return teacherList, result.Error
}

func GetTotalPageNum(name string) (*int64, error) {
	var pageNum int64
	result := database.DB.Model(models.Teacher{}).Where(models.User{
		Username: name,
	}).Count(&pageNum)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pageNum, nil
}

func GetTeacherByTeacherID(teacherID int) (*models.Teacher, int, error) {
	var teacher *models.Teacher
	result := database.DB.Preload("Students").Where("id = ?", teacherID).First(&teacher)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	if teacher.Students == nil {
		return teacher, 0, nil
	}
	return teacher, len(teacher.Students), nil
}

func GetTeacherDDLByUserID(userId int) (models.DDL, error) {
	var ddl models.DDL
	result := database.DB.Where(models.DDL{
		UserID:  userId,
		DDLType: 1,
	}).First(&ddl)
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

func UpdateSelectionTable(userId int, document string) error {
	var student *models.Student
	if err := database.DB.Where(models.Student{
		UserID: userId,
	}).First(&student).Error; err != nil {
		return err
	}
	student.SelectionTable = document
	student.AdminStatus = 1
	if err := database.DB.Omit("teacher_id").Save(&student).Error; err != nil {
		return err
	}
	return nil
}

func aseEncryptStudentInfo(student *models.Student) {
	student.Email = utils.AesEncrypt(student.Email)
	student.Phone = utils.AesEncrypt(student.Phone)
	student.Address = utils.AesEncrypt(student.Address)
	student.Plan = utils.AesEncrypt(student.Plan)
	student.Experience = utils.AesEncrypt(student.Experience)
	student.Honor = utils.AesEncrypt(student.Honor)
	student.Interest = utils.AesEncrypt(student.Interest)
}

func aseDecryptStudentInfo(student *models.Student) {
	student.Email = utils.AesDecrypt(student.Email)
	student.Phone = utils.AesDecrypt(student.Phone)
	student.Address = utils.AesDecrypt(student.Address)
	student.Plan = utils.AesDecrypt(student.Plan)
	student.Experience = utils.AesDecrypt(student.Experience)
	student.Honor = utils.AesDecrypt(student.Honor)
	student.Interest = utils.AesDecrypt(student.Interest)
}

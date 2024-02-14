package teacherService

import (
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/utils"
	"SelectionSystem-Back/config/database"
	"gorm.io/gorm"
	"time"
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
	result := database.DB.Model(models.Student{}).Where(&models.Student{
		TargetID: targetId,
	}).Where("target_status = ? AND admin_status = ?", 0, 0).Find(&studentList)
	for i := range studentList {
		AseDecryptStudentInfo(&studentList[i])
	}
	return studentList, result.Error
}

func StudentCheckList(targetId int) ([]models.Student, error) {
	var studentList []models.Student
	result := database.DB.Model(models.Student{}).Where(models.Student{
		TargetID: targetId,
	}).Where("target_status IN (?)", []int{1, 2}).Find(&studentList)
	for i := range studentList {
		AseDecryptStudentInfo(&studentList[i])
	}
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
	AseDecryptStudentInfo(&info)
	return &info, nil
}

func GetStudentInfoByStudentID(studentId string) (*models.Student, error) {
	var info models.Student
	result := database.DB.Where(&models.Student{
		StudentID: studentId,
	}).First(&info)
	if result.Error == gorm.ErrRecordNotFound {
		info.StudentID = studentId
		return &info, result.Error
	} else if result.Error != nil {
		return nil, result.Error
	}
	AseDecryptStudentInfo(&info)
	return &info, nil
}

func GetTeacherByTeacherID(teacherId int) (*models.Teacher, error) {
	var info models.Teacher
	result := database.DB.Where(&models.Teacher{
		ID: teacherId,
	}).First(&info)
	if result.Error == gorm.ErrRecordNotFound {
		info.ID = teacherId
		return &info, result.Error
	} else if result.Error != nil {
		return nil, result.Error
	}
	return &info, nil
}

func GetTeacherByUserID(userID int) (*models.Teacher, int, error) {
	var teacher *models.Teacher
	result := database.DB.Preload("Students").Where("user_id = ?", userID).First(&teacher)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	if teacher.Students == nil {
		return teacher, 0, nil
	}
	for i := range teacher.Students {
		AseDecryptStudentInfo(&teacher.Students[i])
	}
	return teacher, len(teacher.Students), nil
}

func UpdateStudentInfo(Id int, info *models.Student) error {
	AseEncryptStudentInfo(info)
	result := database.DB.Model(models.Student{}).Where(&models.Student{
		ID: Id,
	}).Updates(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateStudentInfoByStudentID(studentId string, info *models.Student) error {
	AseEncryptStudentInfo(info)
	result := database.DB.Model(models.Student{}).Where(&models.Student{
		StudentID: studentId,
	}).Updates(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func SetDDL(time time.Time, check, userId int) error {
	var result *gorm.DB
	if check == 1 {
		result = database.DB.Model(&models.DDL{}).Where(models.DDL{UserID: userId, DDLType: 1}).Update("first_ddl", time)
		return result.Error
	} else if check == 2 {
		result = database.DB.Model(&models.DDL{}).Where(models.DDL{UserID: userId, DDLType: 1}).Update("second_ddl", time)
		return result.Error
	}
	return result.Error
}

func GetAdminDDL() (models.DDL, error) {
	var ddl models.DDL
	result := database.DB.Where(models.DDL{DDLType: 2}).First(&ddl)
	return ddl, result.Error
}

func UpdateTeacher(teacher *models.Teacher) error {
	result := database.DB.Model(models.Teacher{}).Updates(&teacher)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetStudentList(teacherId int) ([]models.Student, error) {
	var studentList []models.Student
	result := database.DB.Model(models.Student{}).Where(models.Student{
		TeacherID: teacherId,
	}).Find(&studentList)
	for i := range studentList {
		AseDecryptStudentInfo(&studentList[i])
	}
	return studentList, result.Error
}

func AseEncryptStudentInfo(student *models.Student) {
	student.Email = utils.AesEncrypt(student.Email)
	student.Phone = utils.AesEncrypt(student.Phone)
	student.Address = utils.AesEncrypt(student.Address)
	student.Plan = utils.AesEncrypt(student.Plan)
	student.Experience = utils.AesEncrypt(student.Experience)
	student.Honor = utils.AesEncrypt(student.Honor)
	student.Interest = utils.AesEncrypt(student.Interest)
}

func AseDecryptStudentInfo(student *models.Student) {
	student.Email = utils.AesDecrypt(student.Email)
	student.Phone = utils.AesDecrypt(student.Phone)
	student.Address = utils.AesDecrypt(student.Address)
	student.Plan = utils.AesDecrypt(student.Plan)
	student.Experience = utils.AesDecrypt(student.Experience)
	student.Honor = utils.AesDecrypt(student.Honor)
	student.Interest = utils.AesDecrypt(student.Interest)
}

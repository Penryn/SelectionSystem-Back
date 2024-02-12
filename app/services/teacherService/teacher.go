package teacherService

import (
	"SelectionSystem-Back/app/models"
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

func GetTeacherByUserID(userId int) (*models.Teacher, error) {
	var info models.Teacher
	result := database.DB.Where(&models.Teacher{
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

func UpdateStudentInfo(Id int, info *models.Student) error {
	result := database.DB.Model(models.Student{}).Where(&models.Student{
		ID: Id,
	}).Updates(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateStudentInfoByStudentID(studentId string, info *models.Student) error {
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

func UpdateTeacher(Id int, teacher *models.Teacher) error {
	result := database.DB.Model(models.Teacher{}).Where(&models.Teacher{
		ID: Id,
	}).Updates(&teacher)
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
	return studentList, result.Error
}

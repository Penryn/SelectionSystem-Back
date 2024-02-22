package teacherService

import (
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/services/userService"
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

func StudentList(targetId int, checkStudentList int) ([]models.Student, error) {
	var studentList []models.Student
	if checkStudentList == 1 {
		result := database.DB.Model(models.Student{}).Where(&models.Student{
			TargetID:     targetId,
			TargetStatus: 1,
		}).Find(&studentList)
		if len(studentList) == 0 {
			return []models.Student{}, nil
		}
		for i := range studentList {
			aseDecryptStudentInfo(&studentList[i])
		}
		return studentList, result.Error
	} else if checkStudentList == 2 {
		result := database.DB.Model(models.Student{}).
			Where("target_id = ? AND target_status IN (?)", targetId, []int{2, 3}).
			Find(&studentList)
		if len(studentList) == 0 {
			return []models.Student{}, nil
		}
		for i := range studentList {
			aseDecryptStudentInfo(&studentList[i])
		}
		return studentList, result.Error
	} else {
		return []models.Student{}, nil
	}
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
	aseDecryptStudentInfo(&info)
	return &info, nil
}

func CheckStudent(studentId string, targetId int) (*models.Student, error) {
	var info models.Student
	result := database.DB.Where(&models.Student{
		StudentID: studentId,
		TargetID:  targetId,
	}).First(&info)
	if result.Error == gorm.ErrRecordNotFound {
		info.StudentID = studentId
		return &info, result.Error
	} else if result.Error != nil {
		return nil, result.Error
	}
	aseDecryptStudentInfo(&info)
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
	return teacher, len(teacher.Students), nil
}

func GetTeacherByID(id int) (*models.Teacher, int, error) {
	var teacher *models.Teacher
	result := database.DB.Preload("Students").Where("id = ?", id).First(&teacher)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	if teacher.Students == nil {
		return teacher, 0, nil
	}
	return teacher, len(teacher.Students), nil
}

func GetStudentsNumByTarget(targetId int) (int64, error) {
	var studentsNum int64
	result := database.DB.Model(models.Student{}).Where("target_id = ? AND target_status = ?", targetId, 1).Count(&studentsNum)
	if result.Error != nil {
		return 0, result.Error
	}
	return studentsNum, nil
}

func GetStudentsByUserID(userID int) ([]models.Student, int, error) {
	var teacher *models.Teacher
	result := database.DB.Preload("Students").Where("user_id = ?", userID).First(&teacher)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	if teacher.Students == nil {
		return []models.Student{}, 0, nil
	}
	students := teacher.Students
	for i := range teacher.Students {
		aseDecryptStudentInfo(&teacher.Students[i])
	}
	return students, len(students), nil
}

func UpdateStudentInfo(studentId string, info *models.Student) error {
	aseEncryptStudentInfo(info)
	result := database.DB.Model(models.Student{}).Where(&models.Student{
		StudentID: studentId,
	}).Updates(&info)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateStudentInfoTargetStatus(studentId string, targetStatus int) error {
	result := database.DB.Model(models.Student{}).Where("student_id", studentId).Update("target_status", targetStatus)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func SetDDL(time time.Time, userId int) error {
	var result *gorm.DB
	result = database.DB.Model(&models.DDL{}).Where(models.DDL{UserID: userId, DDLType: 1}).Update("first_ddl", time)
	return result.Error
}

func GetAdminDDL() (models.DDL, error) {
	var ddl models.DDL
	result := database.DB.Where(models.DDL{DDLType: 2}).First(&ddl)
	return ddl, result.Error
}

func UpdateTeacher(id int, studentsNum int) error {
	result := database.DB.Model(&models.Teacher{}).Where("user_id = ?", id).UpdateColumn("students_num", studentsNum)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func StudentJoinTeacher(studentID string, targetId int) error {
	var student models.Student
	database.DB.Take(&student, "student_id = ?", studentID)
	var teacher models.Teacher
	database.DB.Take(&teacher, "id = ?", targetId)
	student, err := userService.GetStudentByID(student.UserID)
	if err != nil {
		return err
	}
	err = database.DB.Model(&teacher).Association("Students").Append(&student)
	return err
}

func Disassociate(studentID string, targetId int) error {
	var student models.Student
	database.DB.Take(&student, "student_id = ?", studentID)
	var teacher models.Teacher
	database.DB.Take(&teacher, "id = ?", targetId)
	student, err := userService.GetStudentByID(student.UserID)
	if err != nil {
		return err
	}
	err = database.DB.Model(&teacher).Association("Students").Delete(&student)
	if err != nil {
		return err
	}
	result := database.DB.Model(&student).Updates(map[string]interface{}{"target_status": 0, "target_id": 0})
	if result.Error != nil {
		return result.Error
	}
	return err
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

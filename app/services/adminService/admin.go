package adminService

import (
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/app/utils"
	"SelectionSystem-Back/config/database"
	"time"

	"gorm.io/gorm"
)

func SetDDL(time time.Time, ddltype, id int) error {
	var result *gorm.DB
	if ddltype == 1 {
		result = database.DB.Model(&models.DDL{}).Where(models.DDL{UserID: id, DDLType: 2}).Update("first_ddl", time)
		return result.Error
	} else if ddltype == 2 {
		result = database.DB.Model(&models.DDL{}).Where(models.DDL{UserID: id, DDLType: 2}).Update("second_ddl", time)
		return result.Error
	}
	return result.Error
}

func GetAdvices(pagenum, pagesize int) ([]models.Advice, *int64, error) {
	var advices []models.Advice
	var num int64
	result := database.DB.Model(&models.Advice{}).Count(&num)
	if result.Error != nil {
		return advices, nil, result.Error
	}
	result = database.DB.Limit(pagesize).Offset((pagenum - 1) * pagesize).Find(&advices)
	return advices, &num, result.Error
}

func GetUsers(pagenum, pagesize int, name string) ([]models.User, *int64, error) {
	var users []models.User
	var num int64
	query := database.DB.Model(&models.User{})
	if name != "" {
		query = query.Where("username LIKE ?", "%"+name+"%")
	}
	result := query.Count(&num)
	if result.Error != nil {
		return users, nil, result.Error
	}
	result = query.Limit(pagesize).Offset((pagenum - 1) * pagesize).Find(&users)
	//解密
	for i := 0; i < len(users); i++ {
		users[i].Password = utils.AesDecrypt(users[i].Password)
	}
	return users, &num, result.Error
}

func GetStudents(pagenum int, pagesize int) ([]models.Student, *int64, error) {
	var students []models.Student
	var num int64
	result := database.DB.Model(&models.Student{}).Not("selection_table=?", "").Count(&num)
	if result.Error != nil {
		return students, nil, result.Error
	}
	result = database.DB.Not("selection_table=?", "").Find(&students)
	return students, &num, result.Error
}

func ResetPassword(user_id int) error {
	password := aseEncrypt("123456")
	result := database.DB.Model(&models.User{ID: user_id}).Update("password", password)
	return result.Error
}

func aseEncrypt(data string) string {
	return utils.AesEncrypt(data)
}

func CheckTable(studentID string, target_id int, check int) error {
	var student models.Student
	result := database.DB.Where("student_id = ?", studentID).First(&student)
	if result.Error != nil {
		return result.Error
	}
	if check == 1 {
		err := StudentjoinTeacher(studentID, target_id)
		if err != nil {
			return err
		}
		result = database.DB.Model(&student).Update("admin_status", 2)
	} else if check == 2 {
		result = database.DB.Model(&student).Update("admin_status", 3)
	}
	return result.Error

}

func StudentjoinTeacher(studentID string, target_id int) error {
	var student models.Student
	database.DB.Take(&student, "student_id = ?", studentID)
	var teacher models.Teacher
	database.DB.Take(&teacher, "id = ?", target_id)
	student, err := userService.GetStudentByID(student.UserID)
	if err != nil {
		return err
	}
	err = database.DB.Model(&teacher).Association("Students").Append(&student)
	return err
}

func GetCheckStudents(check int, name string, studentid string) ([]models.Student, error) {
	var students []models.Student
	var result *gorm.DB
	if check == 1 {
		result = database.DB.Where(models.Student{
			TargetStatus: 2,
			AdminStatus: 1,
			Name:        name,
			StudentID:   studentid,
		}).Find(&students)
	} else if check == 2 {
		result = database.DB.Where(models.Student{
			TargetStatus: 2,
			AdminStatus: 2,
			Name:        name,
			StudentID:   studentid,
		}).Or(models.Student{
			TargetStatus: 2,
			AdminStatus: 3,
			Name:        name,
			StudentID:   studentid,
		}).Find(&students)
	}
	return students, result.Error
}

func Disassociate(studentID string, target_id int) error {
	var student models.Student
	database.DB.Take(&student, "student_id = ?", studentID)
	var teacher models.Teacher
	database.DB.Take(&teacher, "id = ?", target_id)
	student, err := userService.GetStudentByID(student.UserID)
	if err != nil {
		return err
	}
	var result *gorm.DB 
	if student.AdminStatus == 2 {
		err = database.DB.Model(&teacher).Association("Students").Delete(&student)
		if err != nil {
			return err
		}
		result = database.DB.Model(&student).Updates(map[string]interface{}{"admin_status": 1})
		err =userService.SendConversation(1, student.UserID, "您被管理员通过的双向选择已被取消")
		if err != nil {
			return err
		}
	}else if student.AdminStatus == 3 {
		result = database.DB.Model(&student).Updates(map[string]interface{}{"admin_status": 1})
		err =userService.SendConversation(1, student.UserID, "您被管理员驳回的双向选择已被取消")
		if err != nil {
			return err
		}
	}
	return result.Error
}

func GetTeachers(pagenum, pagesize int) ([]models.Teacher, *int64, error) {
	var num int64
	var teachers []models.Teacher
	result := database.DB.Model(&models.Teacher{}).Count(&num)
	if result.Error != nil {
		return teachers, nil, result.Error
	}
	result = database.DB.Preload("Students").Limit(pagesize).Offset((pagenum - 1) * pagesize).Find(&teachers)
	return teachers, &num, result.Error
}

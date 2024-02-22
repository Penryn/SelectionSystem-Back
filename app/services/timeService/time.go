package timeService

import (
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/config/database"
	"math/rand"
	"time"
)

func QueryStudents() ([]models.Student, error) {
	var students []models.Student
	result := database.DB.Find(&students)
	return students, result.Error
}

func QueryTeachers() ([]int, error) {
	var teachers []models.Teacher
	result := database.DB.Where("students_num < ?",6).Find(&teachers)
	if result.Error != nil {
		return nil, result.Error
	}
	var teachers_id []int
	for _, teacher := range teachers {
		teachers_id = append(teachers_id, teacher.ID)
	}
	return teachers_id, nil
}

func QueryTime() (time.Time, error) {
	var ddl models.DDL
	result := database.DB.Where("ddl_type = ?",2).First(&ddl)
	return ddl.FirstDDL, result.Error
}


func RandomTeacher(teachers_id []int) int {
	if len(teachers_id) == 0 {
		return 0
	}
	index := rand.Intn(len(teachers_id))
	return teachers_id[index]
}

func UpdateStudent(student models.Student) error {
	result := database.DB.Model(&models.Student{}).Where("id=?", student.ID).Updates(map[string]interface{}{"target_id": student.TargetID, "target_status": 2})
	return result.Error
}

func UpdateTeacher(teacher_id int) error {
	result := database.DB.Model(&models.Teacher{}).Where("id=?", teacher_id).Update("students_num", database.DB.Raw("students_num+1"))
	return result.Error
}
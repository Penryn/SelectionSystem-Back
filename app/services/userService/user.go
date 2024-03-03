package userService

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"

	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/utils"
	"SelectionSystem-Back/config/config"
	"SelectionSystem-Back/config/database"
	"math/rand"
)

func CreateUser(user models.User) error {
	AseEncryptPassword(&user)
	result_a := database.DB.Create(&user)
	var student models.Student
	student.UserID = user.ID
	student.StudentID = user.Username
	student.Name = "未填写"
	student.Email = "未填写"
	student.Class = "未填写"
	student.Phone = "未填写"
	student.Address = "未填写"
	student.PoliticalStatus = "未填写"
	student.Plan = "未填写"
	student.Experience = "未填写"
	student.Honor = "未填写"
	student.Interest = "未填写"
	AseEncryptStudentInfo(&student)
	result_b := database.DB.Omit("teacher_id").Create(&student)
	if result_a.Error != nil {
		return result_a.Error
	} else if result_b.Error != nil {
		return result_b.Error
	} else {
		return nil
	}
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	result := database.DB.Where(models.User{Username: username}).First(&user)
	if user.Password != "" {
		AseDecryptPassword(&user)
	}
	return user, result.Error
}

func GetAvartar() string {
	avatars := []string{
		"https://phlin.love/static/5e75bf70-ce05-4920-b7ba-7b8bcf1c9581.jpg",
		"https://phlin.love/static/712df083-8df1-40a4-bdea-e6403b25314a.jpg",
		"https://phlin.love/static/be262a7c-fe00-4108-bcd4-b0609e9baca5.jpg",
		"https://phlin.love/static/7433cca3-86c1-4ea3-8e24-840f53ac85b8.jpg",
		"https://phlin.love/static/cacc8913-f6c0-45ad-8b84-c3871c664f94.jpg",
		"https://phlin.love/static/e515b3f9-d7f1-4e40-9086-8ded75d0838b.jpg",
		"https://phlin.love/static/c86ae912-ca99-4637-9228-ece8fa06aedf.jpg",
		"https://phlin.love/static/42c06eca-7bf9-4d42-819c-b46295be80b1.jpg",
		"https://phlin.love/static/42c06eca-7bf9-4d42-819c-b46295be80b1.jpg",
		"https://phlin.love/static/4d138da9-be71-43a7-8018-3a90b6ef5f9f.jpg",
		"https://phlin.love/static/5bbe2abb-58a5-447c-80e1-e828b84d330f.jpg",
		"https://phlin.love/static/3f11138a-7a9e-4aaa-ae48-77358590ba42.jpg",
		"https://phlin.love/static/dc07b05a-2711-492a-9962-1d289ef61d7d.jpg",




		
	}
	randomIndex := rand.Intn(len(avatars))
	return avatars[randomIndex]
}

func GetUserByID(id int) (models.User, error) {
	var user models.User
	result := database.DB.Where(models.User{ID: id}).First(&user)
	AseDecryptPassword(&user)
	return user, result.Error
}

func GetStudentByStudentID(sid string) (models.Student, error) {
	var student models.Student
	result := database.DB.Where(models.Student{StudentID: sid}).First(&student)
	AseDecryptStudentInfo(&student)
	return student, result.Error
}

func GetStudentByStudentIDAndAdminStatus(sid string) (models.Student, error) {
	var student models.Student
	result := database.DB.Where(models.Student{StudentID: sid}).Not(models.Student{AdminStatus: 0}).First(&student)
	AseDecryptStudentInfo(&student)
	return student, result.Error
}

func GetStudentByID(id int) (models.Student, error) {
	var student models.Student
	result := database.DB.Where(models.Student{UserID: id}).First(&student)
	AseDecryptStudentInfo(&student)
	return student, result.Error
}


func UpdatePassword(user models.User, newPassword string) error {
	password :=utils.AesEncrypt(newPassword)
	result := database.DB.Model(&user).Update("password", password)
	return result.Error
}

// 新建管理员
func CreateAdministrator() error {
	uname := config.Config.GetString("Administrator.Name")
	upass := config.Config.GetString("Administrator.Pass")
	_, err := GetUserByUsername(uname)
	if err == nil {
		return nil
	}
	var user models.User
	user.Username = uname
	user.Password = upass
	user.Type = 3
	user.Avartar = GetAvartar()
	AseEncryptPassword(&user)
	result := database.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	user, err = GetUserByUsername(uname)
	if err != nil {
		return err
	}
	result = database.DB.Create(&models.DDL{UserID: user.ID, DDLType: 2, FirstDDL: time.Now().AddDate(0, 1, 0), SecondDDL: time.Now().AddDate(0,2,0)})
	return result.Error
}

// 导入教师excel表
func ImportTeacherExcel() error {
	file, err := excelize.OpenFile("德育导师名单.xlsx")
	if err != nil {
		return err
	}
	records, err := file.GetRows("Sheet1")
	if err != nil {
		return err
	}
	for i, record := range records {
		if i == 0 || i == 1 {
			continue
		}
		username, err := strconv.Atoi(record[0])
		if err != nil {
			return err
		}
		result := database.DB.Where(models.User{Username: fmt.Sprintf("zjut%03d", username)}).First(&models.User{})
		if result.Error == gorm.ErrRecordNotFound {
			var user models.User
			user.Username = fmt.Sprintf("zjut%03d", username)
			user.Password = "123456"
			user.Type = 2
			user.Avartar = GetAvartar()
			AseEncryptPassword(&user)
			result := database.DB.Create(&user)
			if result.Error != nil {
				return result.Error
			}
			user, err := GetUserByUsername(fmt.Sprintf("zjut%03d", username))
			if err != nil {
				return err
			}
			result = database.DB.Create(&models.Teacher{
				UserID:      user.ID,
				TeacherName: record[1],
				Section:     record[2],
				Office:      record[3],
				Phone:       record[4],
				Email:       record[5],
			})
			if result.Error != nil {
				return result.Error
			}
			result = database.DB.Omit("second_ddl").Create(&models.DDL{UserID: user.ID, DDLType: 1, FirstDDL: time.Now().AddDate(0, 1, 0)})
			if result.Error != nil {
				return result.Error
			}
		} else if result.Error == nil {
			return nil
		} else {
			return result.Error
		}
	}
	return nil
}

func SendConversation(userAID int, userBID int, message string) error {
	result := database.DB.Create(&models.Conversation{UserAID: userAID, UserBID: userBID, Content: message, Time: time.Now()})
	return result.Error
}

func GetConversation(userAID int, userBID int) ([]models.Conversation, error) {
	var conversation []models.Conversation
	result := database.DB.Where(models.Conversation{UserAID: userAID, UserBID: userBID}).Or(models.Conversation{UserAID: userBID, UserBID: userAID}).Find(&conversation)
	return conversation, result.Error
}

func CreateReason(userID int, reason_name, reason_content string) error {
	result := database.DB.Create(&models.Reason{UserID: userID, ReasonName: reason_name, ReasonContent: reason_content})
	return result.Error
}

func UpdateReason(userID int, reasonID int, reason_name, reason_content string) error {
	result := database.DB.Model(&models.Reason{ID: reasonID, UserID: userID}).Updates(models.Reason{ReasonName: reason_name, ReasonContent: reason_content})
	return result.Error
}

func DeleteReason(userID int, reasonID int) error {
	result := database.DB.Where(models.Reason{ID: reasonID, UserID: userID}).Delete(&models.Reason{})
	return result.Error
}

func GetReasons(userID int) ([]models.Reason, error) {
	var reasons []models.Reason
	result := database.DB.Where(models.Reason{UserID: userID}).Find(&reasons)
	return reasons, result.Error
}

func GetReasonsByReasonID(reasonID int) (models.Reason, error) {
	var reason models.Reason
	result := database.DB.Where(models.Reason{ID: reasonID}).First(&reason)
	return reason, result.Error
}

func GetReasonByID(id int) (models.Reason, error) {
	var reason models.Reason
	result := database.DB.Where(models.Reason{ID: id}).First(&reason)
	return reason, result.Error
}

func PostReason(userAID int, userBID int, message string) error {
	result := database.DB.Create(&models.Conversation{UserAID: userAID, UserBID: userBID, Content: message, Time: time.Now()})
	return result.Error
}


func GetAdminDDL() (models.DDL,error) {
	var ddl models.DDL
	result:=database.DB.Where(models.DDL{DDLType: 2}).First(&ddl)
	return ddl,result.Error
}

func GetTeacherByTeacherID(teacherID int) (models.Teacher,error) {
	var teacher models.Teacher
	result:=database.DB.Where(models.Teacher{UserID: teacherID}).First(&teacher)
	return teacher,result.Error
}

func GetTeacherDDLTime(userID int) (models.DDL,error) {
	var ddl models.DDL
	result:=database.DB.Where(models.DDL{UserID: userID,DDLType: 1}).First(&ddl)
	return ddl,result.Error
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

func GetStudentByUserID(userID int) (models.Student,error) {
	var student models.Student
	result:=database.DB.Where(models.Student{UserID: userID}).First(&student)
	AseDecryptStudentInfo(&student)
	return student,result.Error
}

func GetTeacherByUserID(userID int) (models.Teacher,error) {
	var teacher models.Teacher
	result:=database.DB.Where(models.Teacher{UserID: userID}).First(&teacher)
	return teacher,result.Error
}
func GetReasonByName(reasonName string) (models.Reason,error) {
	var reason models.Reason
	result:=database.DB.Where(models.Reason{ReasonName: reasonName}).First(&reason)
	return reason,result.Error
}

func AseEncryptPassword(user *models.User){
	user.Password = utils.AesEncrypt(user.Password)
}

func AseDecryptPassword(user *models.User){
	user.Password = utils.AesDecrypt(user.Password)
}

func AseEncryptStudentInfo(student *models.Student){
	student.Email = utils.AesEncrypt(student.Email)
	student.Phone = utils.AesEncrypt(student.Phone)
	student.Address = utils.AesEncrypt(student.Address)
	student.Plan = utils.AesEncrypt(student.Plan)
	student.Experience = utils.AesEncrypt(student.Experience)
	student.Honor = utils.AesEncrypt(student.Honor)
	student.Interest = utils.AesEncrypt(student.Interest)
}

func AseDecryptStudentInfo(student *models.Student){
	student.Email = utils.AesDecrypt(student.Email)
	student.Phone = utils.AesDecrypt(student.Phone)
	student.Address = utils.AesDecrypt(student.Address)
	student.Plan = utils.AesDecrypt(student.Plan)
	student.Experience = utils.AesDecrypt(student.Experience)
	student.Honor = utils.AesDecrypt(student.Honor)
	student.Interest = utils.AesDecrypt(student.Interest)
}


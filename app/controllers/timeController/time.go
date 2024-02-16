package timeController

import (
	"SelectionSystem-Back/app/models"
	"SelectionSystem-Back/app/services/timeService"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

var allocated bool

func Time() {
	t := cron.New()
	t.Start()
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		last_time := time.Now()
		for {
			timeValue, err := timeService.QueryTime()
			if err != nil {
				log.Fatal("Server start error:", err)
				return
			} else {
				if timeValue != last_time {
					if timeValue.Before(time.Now()) {
						timeValue = time.Now()
					}
					last_time = timeValue
					t.Stop()
					duringtime := time.Until(timeValue)
					sh := cron.Every(duringtime)
					fmt.Println(sh)
					t.Schedule(sh, cron.FuncJob(func() {
						if !allocated {
							Allocate()
							allocated = true
						}
					}))

					t.Start()
				}
			}
			time.Sleep(5 * time.Minute)
		}
	}()

	// 等待关闭信号
	<-stopCh

	// 优雅关闭定时器
	t.Stop()
}

func Allocate() {
	// 获取所有学生
	var students []models.Student
	students, err := timeService.QueryStudents()
	if err != nil {
		log.Fatal("Server start error:", err)
		return
	}
	fmt.Println(123)
	for _, student := range students {
		if student.TargetID == 0 {
			//获取所有教师
			teachers_id, err := timeService.QueryTeachers()
			if err != nil {
				log.Fatal("Server start error:", err)
				return
			}
			//分配教师
			student.TargetID = timeService.RandomTeacher(teachers_id)
			//更新学生信息
			err = timeService.UpdateStudent(student)
			if err != nil {
				log.Fatal("Server start error:", err)
				return
			}
			//更新教师信息
			err = timeService.UpdateTeacher(student.TargetID)
			if err != nil {
				log.Fatal("Server start error:", err)
				return
			}

		}
	}
}

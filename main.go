package main

import (
	"SelectionSystem-Back/app/controllers/timeController"
	"SelectionSystem-Back/app/midwares"
	"SelectionSystem-Back/app/services/userService"
	"SelectionSystem-Back/config/database"
	"SelectionSystem-Back/config/router"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(midwares.Corss)
	r.Use(midwares.RateLimitMiddleware(time.Second,100,100))
	r.Use(midwares.ErrHandler())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)
	userService.CreateAdministrator()
	userService.ImportTeacherExcel()
	r.Static("/static","./static")
	go timeController.Time()
	router.Init(r)
	err := r.Run()
	if err != nil {
		log.Fatal("Server shutdown error:", err)
	}
}

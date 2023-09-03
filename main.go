package main

import (
	"os"
	"task-5-pbi-btpns-daffajatmiko/database"
	"task-5-pbi-btpns-daffajatmiko/models"
	"task-5-pbi-btpns-daffajatmiko/router"
)

func main() {
	db := database.ConnectDB()
	db.AutoMigrate(&models.User{})

	r := router.InitRoutes(db)
	r.Run(":" + os.Getenv("PORT"))
}
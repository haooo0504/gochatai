package main

import (
	"gochatai/models"
	"gochatai/router"
	"gochatai/utils"
	"log"
)

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Description for what is this security definition being used
func main() {
	utils.InitConfig()
	utils.InitMySQL()

	// Add AutoMigrate calls here
	// 進行數據庫遷移
	db := utils.GetDB()
	err := db.AutoMigrate(&models.UserBasic{})
	if err != nil {
		log.Fatalf("failed to auto migrate UserBasic: %v", err)
	}
	err = db.AutoMigrate(&models.PostInfo{})
	if err != nil {
		log.Fatalf("failed to auto migrate PostInfo: %v", err)
	}

	r := router.Router()
	r.Run(":8082")
}

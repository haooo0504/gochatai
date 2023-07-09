package main

import (
	"gochatai/router"
	"gochatai/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySQL()
	r := router.Router()
	r.Run(":8082")
}

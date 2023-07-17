package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func GetDB() *gorm.DB {
	return DB
}

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config") // 請注意路徑 /config | config
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app:", viper.Get("app"))
	fmt.Println("config app:", viper.Get("mysql"))
}

func InitMySQL() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	var err error
	for i := 0; i < 10; i++ { // Try to connect 10 times
		DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
		if err != nil {
			fmt.Printf((viper.GetString("mysql.dns")))

			// fmt.Printf("failed to connect to the database: %v. Retrying in 10 seconds...\n", err)
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// user := models.UserBasic{}
	// DB.Find(&user)
	// fmt.Println(user, 111)

}

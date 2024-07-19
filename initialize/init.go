package initialize

import (
	"fmt"
	"log"
	"newbee/global"

	// "newbee/models/manage"

	// "newbee/models/mall"
	"os"
	"time"

	// "github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)


func InitConfig() {
	viper.SetConfigName("newbee.env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config mysql:", viper.Get("mysql"))
}

func Init_DB() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	var err error
	global.DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	// DB.First()
	// fmt.Println(DB.Statement)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(" MySQL inited.......")
	// global.DB.AutoMigrate(&mall.MallUserToken{})
	// global.DB.AutoMigrate(&mall.MallUser{})
	// global.DB.AutoMigrate(&mall.MallUserAddress{})
	// global.DB.AutoMigrate(&mall.MallShoppingCartItem{})
	// global.DB.AutoMigrate(&mall.MallMessage{})
	// global.DB.AutoMigrate(&manage.MallAdminUser{})
	// global.DB.AutoMigrate(&manage.MallAdminUserToken{})
	// global.DB.AutoMigrate(&manage.MallGoodsCategory{})
	// global.DB.AutoMigrate(&manage.MallGoodsInfo{})
	// global.DB.AutoMigrate(&manage.MallOrder{})
	// global.DB.AutoMigrate(&manage.MallOrderItem{})
	// global.DB.AutoMigrate(&manage.MallIndexConfig{})
}

// func Init_Redis() {
// 	global.Redis = redis.NewClient(&redis.Options{
// 		Addr:         viper.GetString("redis.addr"),
// 		Password:     viper.GetString("redis.password"),
// 		DB:           viper.GetInt("redis.DB"),
// 		PoolSize:     viper.GetInt("redis.poolSize"),
// 		MinIdleConns: viper.GetInt("redis.minIdleConn"),
// 	})
// }
func Init() {
	InitConfig()
	Init_DB()
	// Init_Redis()
}
// 创建时间：2026/8/17 下午6:55
package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"learn_layer_demo/model"
	"log"
)

// DB 全局数据库实例，repository层使用
var DB *gorm.DB

func InitDB() error {
	// 修改为你自己的mysql账号密码、数据库名
	dsn := "root:12345678@tcp(127.0.0.1:3306)/learn_layer_demo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db

	// 新增：自动迁移，根据结构体创建 users 表
	log.Println("开始执行自动迁移...")
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Println("自动迁移失败：", err)
		return err
	}
	log.Println("自动迁移成功，users 表已创建")

	//db.Create(&model.User{Name: "张三", Age: 20})
	//db.Create(&model.User{Name: "李四", Age: 30})
	//db.Create(&model.User{Name: "王五", Age: 40})
	//db.Create(&model.User{Name: "赵六", Age: 50})
	//db.Create(&model.User{Name: "钱七", Age: 60})
	//db.Create(&model.User{Name:"王亮",Age: 20,Username: "宝宝巴士",Password: "abc123"})
	//db.Create(&model.User{Name:"孙超",Age: 54,Username: "马尾卡",Password: "12345678"})
	//db.Create(&model.User{Name:"赵桐",Age: 33,Username: "芙宁娜",Password: "xxx66666"})
	//db.Create(&model.User{Name:"郑飞",Age: 24,Username: "admin",Password: "123456"})

	return nil
}

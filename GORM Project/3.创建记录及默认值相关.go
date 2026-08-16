// 创建时间：2026/8/9 下午8:11
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Man struct {
	ID   uint   `gorm:"primaryKey"`     //标记为主键
	Name string `gorm:"default:'百家姓'"`  //默认值为空字符串，也可以指定默认值为百家姓
	Age  int    `gorm:"column:man_age"` //在数据库中以man_age存在。默认值为0
}

func main() {
	// ========== 1.DSN 数据库连接字符串，必须修改成你自己数据库信息 ==========
	// 格式：用户名:密码@tcp(地址:端口)/数据库名?charset=utf8mb4&parseTime=True&loc=Local
	dsn := "root:12345678@tcp(127.0.0.1:3306)/gorm_db?charset=utf8mb4&parseTime=True&loc=Local"
	// ========== 2.打开连接，获取db对象，后续所有操作都靠db ==========
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 日志配置：logger.Info 打印执行的SQL语句，学习阶段建议开启，方便看到gorm底层生成什么SQL
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("数据库连接失败：", err) // 连接失败直接退出程序
	}
	fmt.Println("✅数据库连接成功")
	// ========== 3.获取底层sqlDB，设置连接池参数（必写，控制连接复用） ==========
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("获取底层db失败：", err)
	}
	sqlDB.SetMaxIdleConns(10)               // 设置最大空闲连接数：空闲保留在池子里的连接
	sqlDB.SetMaxOpenConns(100)              // 设置最大打开连接数：数据库同时允许的最大连接总数
	sqlDB.SetConnMaxLifetime(1 * time.Hour) // 连接最大存活时间：连接超过该时长会被关闭回收
	// ==========4.自动迁移：根据User结构体自动创建/更新数据表（建表） ==========
	// 如果表不存在就创建；字段新增会追加；不会删除原有字段，安全
	err = db.AutoMigrate(&Man{}) //将Man结构体和数据库表man关联
	if err != nil {
		log.Fatal("自动迁移建表失败：", err)
	}
	fmt.Println("✅People表迁移完成")

	//创建记录：数据库中新增一行数据=结构体新增一个实例
	man1 := Man{ID: 0001, Name: "王亮", Age: 25} //指定字段，则使用指定值
	db.Create(&man1)
	man2 := Man{ID: 0002} //未指定字段，则使用默认值，Name: "百家姓", Age: 0
	db.Create(&man2)
	man3 := Man{ID: 0003, Name: "", Age: 30} //即使指定了空值，在数据库中会忽略空值，也会表现为使用默认值。若要避免此情况，可以使用指针或实现Scanner/Valuer接口。示例如下：
	db.Create(&man3)
	/*使用指针
	将结构体定义中改为：Name *string `gorm:"default:'百家姓'"`,并搭配
	man4:=Man{ID: 0004, Name: new(string), Age: 40,}
	*/

	/*使用实现Scanner/Valuer接口
	将结构体定义中改为：Name string sql.NullString `gorm:"default:'百家姓'"`,并搭配
	man5:=Man{ID: 0004, sql.NullString{String: "", Valid: true}, Age: 40,}
	sql.NullString的String字段表示字符串值，Valid字段表示是否为空。当Valid为false时，String字段将被忽略。当Valid为true时，String字段将被使用。
	*/

}

// 创建时间：2026/8/8 下午6:38
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

// User 模型（对应mysql中的 user 表，GORM默认结构体名转小写复数作为表名：users）
type User struct {
	ID   uint   `gorm:"primaryKey"` // 主键，uint无符号整型，自增
	Name string // 对应数据库 name 字段
	Age  int    // 对应数据库 age 字段
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
	// 设置最大空闲连接数：空闲保留在池子里的连接
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数：数据库同时允许的最大连接总数
	sqlDB.SetMaxOpenConns(100)
	// 连接最大存活时间：连接超过该时长会被关闭回收
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	// ==========4.自动迁移：根据User结构体自动创建/更新数据表（建表） ==========
	// 如果表不存在就创建；字段新增会追加；不会删除原有字段，安全
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("自动迁移建表失败：", err)
	}
	fmt.Println("✅表迁移完成")

	// ====================5.【新增 Create】写入数据 =================================
	// 构造结构体实例，准备插入数据库
	u1 := User{ID: 1, Name: "张三", Age: 20}
	// Create(&u1)：传入结构体指针，向users表插入一行记录
	// 注意：传入指针，GORM会回填数据库产生的值（自增ID等）到u1对象
	db.Create(&u1)
	fmt.Println("✅写入数据完成，u1=", u1)

	// ====================6.【查询 First】查询第一条数据 =================================
	var u2 User
	// db.First(&u2)：查询users表第一条记录，结果保存到u2指针指向的变量
	// First 如果找不到记录，gorm会返回ErrRecordNotFound错误
	db.First(&u2)
	fmt.Println("✅查询数据完成，u2=", u2)

	// ====================7.【更新 Update】更新单个字段 =================================
	// Model(&u2)：指定要操作的模型/表，根据u2的主键ID作为where条件
	// Update("name", "李四")：更新name字段为李四；只修改这一个字段
	db.Model(&u2).Update("name", "李四")
	// ⚠️注意：此时内存中的u2变量Name不会自动刷新，还是旧值；需要重新查询才拿到数据库最新数据
	fmt.Println("✅执行更新完成(内存u2不会自动同步数据库), u2=", u2)

	// ====================8.【删除 Delete】软删除 / 删除数据 =================================
	// db.Delete(&u2) 根据u2的主键ID执行删除
	// ⚠️重点：当前模型User没有引入gorm.DeletedAt，执行的是物理删除，数据直接从表消失
	// 如果结构体包含 gorm.DeletedAt，Delete是软删除，只是给记录打上删除标记，数据还在库中
	db.Delete(&u2)
	fmt.Println("✅执行删除完成，u2=", u2)
}

/*
【使用前操作步骤】
1.打开mysql，手动创建数据库：create database gorm_db charset utf8mb4;
2.修改dsn里面的密码为你本机mysql的root密码
3.运行 go run main.go

重要知识点笔记：
1.dsn参数说明
  charset=utf8mb4：支持emoji表情，不要写utf8，mysql的utf8不是真正完整utf‑8
  parseTime=True：必须开启，GORM处理时间类型
  loc=Local：使用本地时区

2.db对象：GORM的核心句柄，后续所有增删改查 db.Create db.Find db.Update db.Delete全部用这个db
⚠️不要在循环、函数内部反复Open打开数据库连接！全局只用一个db实例。

3.AutoMigrate(&User{})自动迁移
  优点：结构体改了，直接运行程序就同步表结构，不用手写建表SQL
  限制：不会删除数据库中已经存在的旧字段，只新增、修改；生产环境谨慎使用，生产一般用迁移脚本。

4.连接池 SetMaxIdleConns / SetMaxOpenConns
  mysql连接是珍贵资源，不能每次操作数据库新建连接；连接池用来复用连接。

5.模型规则
type User struct{} → 默认对应数据库表名为 users（结构体名小写+复数）
ID字段，gorm识别为主键；可以使用标签 gorm:"primaryKey"显式指定主键。

===== 新增CRUD注意点（本次新增注释） =====
① Create(&u1) 务必传指针，GORM会把自增ID回填到结构体对象。
② db.First(&u2)：取表第一条；无数据会返回记录不存在错误。
③ Update更新数据库，但不会自动更新内存中原结构体变量，想要最新数据需要重新查询。
④ Delete物理删除与软删除区分：
   - 没有 gorm.DeletedAt → 物理删除，数据直接删掉
   - 结构体加上 gorm.DeletedAt gorm:"softDelete" → 软删除（逻辑删除）
⑤ Model()方法：用来指定操作的模型，会读取结构体里的主键作为where条件。
*/

// 创建时间：2026/8/15 下午7:58
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Cat struct {
	gorm.Model
	Name   string `gorm:"default:'小黑'"` //默认值为空字符串，也可以指定默认值为小黑
	Age    int
	Weight float64
	Active bool
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
	err = db.AutoMigrate(&Cat{})
	if err != nil {
		log.Fatal("自动迁移建表失败：", err)
	}
	fmt.Println("✅Cat表迁移完成")
	//创建记录：数据库中新增一行数据=结构体新增一个实例。在学习更新时，先创建一些数据，再注释掉。
	//cat1 := Cat{Name: "小黑", Age: 1, Weight: 12.3, Active: true}
	//db.Create(&cat1)
	//cat2 := Cat{Name: "小白", Age: 2, Weight: 16.7, Active: true}
	//db.Create(&cat2)
	//cat3 := Cat{Name: "小黄", Age: 3, Weight: 22.1, Active: true}
	//db.Create(&cat3)

	//更新：Save()默认更新所有字段
	var cat1 Cat    //先创建一个结构体实例，用于接收查询结果，处理更新
	db.First(&cat1) //查询第一条记录，并赋值给cat1
	cat1.Name = "小黑狗"
	cat1.Age = 10
	db.Save(&cat1) //除了Name，Age，其他字段都会被更新 sql：UPDATE `cats` SET `created_at`='2026-08-15 21:04:25.703',`updated_at`='2026-08-15 21:12:48.61',`deleted_at`=NULL,`name`='小黑狗',`age`=10,`weight`=12.3,`active`=true WHERE `cats`.`deleted_at` IS NULL AND `id` = 1

	//更新：使用Update()指定更新单个字段或者Updates指定更新多个字段
	var cat2 Cat
	db.First(&cat2) //查询第一条记录，并赋值给cat2
	//更新单个属性
	db.Model(&cat2).Update("Age", 20) //更新Age字段为20  请注意，update_at会默认更新，不用管
	//更新：根据给定的条件更新单个属性
	db.Model(&cat2).Where("active=?", true).Update("Weight", 100) //当active为true时，更新Weight字段为100
	//更新：：使用map更新多个属性，只会更新其中有变化的属性
	c1 := map[string]interface{}{"Age": 30, "Weight": 200}
	db.Model(&cat2).Updates(c1) //更新Age字段为30，Weight字段为200
	//更新：使用struct更新多个属性，只会更新其中有变化且非零值的属性
	c2 := Cat{Age: 40, Weight: 300, Active: false}
	db.Model(&cat2).Updates(c2) //更新Age字段为40，Weight字段为300，Active字段为false。
	var num int64
	num = db.Model(&cat2).Updates(c2).RowsAffected //RowsAffected能够获取更新影响的行数
	fmt.Println("更新影响行数：", num)

	//上面这些方法，会将c1或c2中的所有内容都更新到数据库中，如果只想更新其中一部分，可以使用Select和Omit方法来指定要更新或忽略的字段。
	var cat3 Cat
	db.First(&cat3) //查询第一条记录，并赋值给cat2
	c3 := Cat{Age: 50, Weight: 400, Active: false}
	db.Model(&cat3).Select("Age,Weight").Updates(c3) //只将c3中的Age和Weight更新到数据库中
	db.Model(&cat3).Omit("Age,Weight").Updates(c3)   //只将c3中的Active更新到数据库中,忽略Age和Weight字段

	//Hooks：上面这些更新操作会自动运行model的BeforeUpdate、AfterUpdate方法，更新UpdateAt时间戳，如果不想调用这些方法，可以使用UpdateColumn、UpdateColumns来更新单个字段或多个字段，不会更新UpdateAt时间戳。另外，执行批量更新Updates时，不会触发BeforeUpdate、AfterUpdate方法。
	var cat4 Cat
	db.First(&cat4) //查询第一条记录，并赋值给cat2
	c4 := Cat{Age: 60, Weight: 500, Active: false}
	db.Model(&cat4).UpdateColumns(c4) //只将c4中的更新到数据库中，不会更新UpdateAt时间戳

	//更新；使用sql表达式更新
	/*
		Model()方法关键区别总结：
		1. db.Model(&Cat{}) 传入【空结构体】Cat{}，结构体主键ID=0
		   - GORM只识别模型对应数据表，**不会自动生成WHERE id条件**
		   - 后续Update/Delete若不手动加Where条件：会执行【全表更新/全表删除】，属于高危操作！
		   - 适用场景：批量操作、Count、Find、Pluck等查询，需要自己写Where做过滤。
		2. db.Model(&cat3) 传入【变量实例】cat3
		   - 情况①：cat3已经经过First/Find查询，内存中主键ID≠0
		     → GORM读取内存里的ID，**自动拼接 WHERE id = 内存中的id**
		     → 后续Update/Delete只操作这一条记录，用于单条更新、单条删除。
		     ⚠️注意：ID取自Go内存变量，不会重新查询数据库。
		   - 情况②：cat3仅var定义、未查询，ID=0（零值）
		     → 和 Model(&Cat{}) 行为一致，不会生成id条件，不加Where会全表操作！
		⚠️高危坑点：
		① 不要传入未查询的零值实例给Model做Update/Delete，会误触发全表更新删除；
		② 使用Model(&空结构体{})做更新删除，必须手动拼接Where条件限定范围；
		③ Model本身不执行SQL，只是组装链式上下文；真正执行数据库操作靠Update、Delete等立即执行方法。
	*/
	db.Model(&Cat{}).Update("age", gorm.Expr("age + ?", 2)) //执行全表更新，表中所有记录的age字段都会加2
	var cat5 Cat
	db.First(&cat5)                                                //查询第一条记录，并赋值给cat2
	db.Model(&cat5).Update("age", gorm.Expr("age * ? + ?", 2, 10)) //执行单条更新，只更新cat5记录的age字段

	//更新：修改Hook中的值
}

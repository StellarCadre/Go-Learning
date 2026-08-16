// 创建时间：2026/8/16 下午5:14
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Dog struct {
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
	err = db.AutoMigrate(&Dog{})
	if err != nil {
		log.Fatal("自动迁移建表失败：", err)
	}
	fmt.Println("✅Cat表迁移完成")
	//创建记录：数据库中新增一行数据=结构体新增一个实例。在学习更新时，先创建一些数据，再注释掉。
	//dog1 := Dog{Name: "小黑", Age: 1, Weight: 12.3, Active: true}
	//db.Create(&dog1)
	//dog2 := Dog{Name: "小白", Age: 2, Weight: 16.7, Active: true}
	//db.Create(&dog2)
	//dog3 := Dog{Name: "小黄", Age: 3, Weight: 22.1, Active: true}
	//db.Create(&dog3)
	//dog4 := Dog{Name: "小花", Age: 4, Weight: 25.5, Active: false}
	//db.Create(&dog4)
	//dog5 := Dog{Name: "小狗", Age: 5, Weight: 28.9, Active: true}
	//db.Create(&dog5)
	//dog6 := Dog{Name: "小猫", Age: 6, Weight: 32.3, Active: true}
	//db.Create(&dog6)

	//删除：
	//错误写法：
	/*不推荐，千万不要这样写
	如果 Find 查询结果为空（没有匹配数据），dog1.ID 会被置为 0。
	接着执行 db.Delete(&dog1)，此时 ID=0，没有 where id 条件 → 全表软删除！高危！
	*/
	//db.Find(&dog1)
	//db.Delete(&dog1)
	/*不推荐，千万不要这样写
	查到数据：dog1.ID为真实主键值，Delete(&dog1)正常只删这一行。
	风险：如果表是空的，First 查不到数据，返回gorm.ErrRecordNotFound，dog1.ID=0，此时再 Delete 依旧会全表删除。
	First / Find 如果查询不到记录，结构体ID变为0；此时直接Delete会发生【全表软删除】高危事故。
	*/
	//db.First(&dog1)
	//fmt.Println(dog1)
	//db.Delete(&dog1)

	//正确写法：
	/*
		GORM软删除正确写法汇总：
		1. 已知主键ID：
		   ① var dog Dog; dog.ID=1; db.Delete(&dog)
		   ② db.Delete(&Dog{}, 1) 简洁写法；批量id传入切片 []int{1,2}
		2. 需要先查询数据再删除：使用First，必须接收并判断error，确认查到记录再Delete，避免ID=0全表删除。
		3. 按业务条件批量删除：必须搭配Where，db.Where(条件).Delete(&Dog{})，禁止不带Where直接Delete(&Dog{})。
		4. Unscoped().Delete() 为物理硬删除，真正从数据库删除该行。

		软删除与硬删除区别：
		1. 软删除：只更新 deleted_at 字段，不删除数据。
		2. 硬删除：直接删除数据。
		软删除不是 DELETE 删除行，是执行 UPDATE，给 deleted_at 字段赋值当前时间，这条记录仍然物理保存在数据表里面。
		普通 GORM 查询会自动过滤掉软删除的数据，当调用 Find / First / Where 普通查询，GORM 会自动拼接条件：WHERE `dogs`.`deleted_at` IS NULL
	*/

	//写法1：你已经知道要删的 id 是多少
	var dog1 Dog //等价写法 var dog1 = Dog{}
	dog1.ID = 6
	db.Delete(&dog1) //等价db.Delete(&Dog{}, 6)
	/*
		警告：删除记录时，请确保主键字段有值，gorm会通过主键去删除对应的记录，如果主键字段为空，则会删除全部记录
	*/
	var dog2 Dog     //初始时全部字段为零值：ID = 0，Name = ""等等
	dog2.Name = "小黑" //这里未指定主键的值，而是指定了其他的
	db.Delete(&dog2) //ID 为零，无法生成 WHERE id = ? 条件，最终 SQL没有任何过滤条件。没有 WHERE id=xxx，表里面所有未被软删除的行全部被打上删除标记，全表软删除！

	//写法2：条件批量软删除，不走结构体的 ID 主键，依靠 Where 条件筛选要删除哪些行
	db.Where("name = ?", "小黑").Delete(&Dog{}) //.Where() 链式设置过滤条件：name = '小黑' Delete(&Dog{})：传入空模型&Dog{}，用来告诉 GORM 操作哪一张表（dogs表）
	db.Where("name like ?", "%小%").Delete(&Dog{})
	db.Delete(&Dog{}, "name like ?", "%白%") //Delete 的内联条件写法，条件直接写在 Delete 的参数里，不用单独写.Where()

	//写法3：硬删除，直接从数据库删除该行
	db.Unscoped().Where("name=?", "小狗").Delete(&Dog{})

}

/*
db.Delete(&dogObj)：传已经实例化的结构体对象 → 取对象内部主键 ID作为删除条件，忽略其他字段。

db.Where(条件).Delete(&Dog{}) / db.Delete(&Dog{}, 条件)：传空模型 &Dog {}
→ 不从结构体拿 ID，完全依靠 Where / 内联条件做批量删除，可以删除一行或者多行。

软删除的本质：
1. gorm.Model内部定义了 deleted_at *time.Time 字段；GORM检测到该字段存在，Delete就执行UPDATE软删除。
2. 如果结构体不嵌入gorm.Model，并且没有手动定义deleted_at字段：
   db.Delete()执行【硬删除】，生成DELETE语句，记录直接从数据表移除。
3. 想要自己实现软删除，必须手动定义字段：deleted_at *time.Time（类型必须是指针time.Time）。
4. 仅仅自己起别的字段名（如delete_time），GORM不会识别为软删除，依旧是硬删除。
*/

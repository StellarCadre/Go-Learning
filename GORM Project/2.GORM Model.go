// 创建时间：2026/8/8 下午8:16
package main

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

/*
gorm.Model 是GORM官方内置的基础结构体（源码如下）
嵌入该结构体后自动获得4个字段：
type Model struct {
    ID        uint           `gorm:"primaryKey"` // 主键ID，自增
	CreatedAt time.Time                          // 创建时间，记录行创建时刻
	UpdatedAt time.Time                          // 更新时间，记录行最后修改时刻
	DeletedAt *time.Time                         // 软删除标记，指针类型，nil代表未删除
}
嵌入 gorm.Model 等价于把这四个字段直接写在自己结构体里面；
并不是真正面向对象的继承，是Go的结构体嵌套（匿名嵌入）。
如果不想使用软删除，可以自己只写ID、CreatedAt、UpdatedAt，不写DeletedAt。
*/

// People 自定义的model模型，匿名嵌入gorm.Model，自动获得ID、CreatedAt、UpdatedAt、DeletedAt
type People struct {
	gorm.Model // 匿名嵌入GORM内置Model，拥有主键、创建、更新、软删除字段

	Name         string        // 普通字符串字段
	Age          sql.NullInt64 // sql.NullInt64：可以为NULL的int类型，区分0值和数据库null
	Birthday     *time.Time    // 指针time.Time，nil代表数据库该字段为NULL
	Email        string        `gorm:"type:varchar(100);uniqueIndex"` // type指定数据库字段类型；uniqueIndex创建唯一索引
	Role         string        `gorm:"size:255"`                      // size：指定字符串字段最大长度
	NumberNumber *string       `gorm:"unique;not null"`               // unique唯一约束；not null非空约束
	Num          int           `gorm:"AUTO_INCREMENT"`                // 设置字段自增
	Address      string        `gorm:"index:addr"`                    // index:addr，给该字段建立名为addr的普通索引
	IgnoreMe     int           `gorm:"-"`                             // gorm:"-" 忽略此字段，不映射到数据库，仅go代码内部使用
}

// Animal 不嵌入gorm.Model，完全自定义模型
type Animal struct {
	AnimalID uint   `gorm:"primaryKey"`      // 自定义主键，不使用默认ID字段
	Name     string `gorm:"column:ani_name"` // column:xxx，手动指定数据库列名；Go结构体字段Name → 数据库列 ani_name
}

// TableName 是GORM约定的方法，自定义表名
// 不需要手动调用！AutoMigrate、CRUD的时候GORM内部会自动调用该函数获取表名
// 默认规则：结构体名转小写复数，如People → peoples；实现TableName()可以覆盖默认表名
func (Animal) TableName() string {
	return "Animalsss" // 指定该模型对应数据库表名叫 Animalsss
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

	// ==========4.自动迁移：根据结构体自动创建/更新数据表（建表） ==========
	// 如果表不存在就创建；字段新增会追加；不会删除原有字段，安全
	err = db.AutoMigrate(&People{})
	if err != nil {
		log.Fatal("自动迁移建表失败：", err)
	}
	fmt.Println("✅People表迁移完成")

	err = db.AutoMigrate(&Animal{})
	if err != nil {
		log.Fatal("自动迁移建表失败：", err)
	}
	fmt.Println("✅Animal表迁移完成")

	/*
		补充：另一种指定表名方式（优先级高于TableName方法）
		db.Table("Animalsss").AutoMigrate(&Animal{})
		直接强制指定操作哪张表；适合临时操作、表名动态变化场景。
	*/
}

/*
==================== GORM模型标签、命名规则总结 ====================
一、gorm.Model内置嵌入
1. 匿名嵌入 gorm.Model，自动带上 ID(主键)、CreatedAt、UpdatedAt、DeletedAt(软删除)
2. DeletedAt不为nil时代表软删除，Delete方法不会真正删除数据，只是给该字段写入删除时间。

二、常用gorm标签（写在结构体字段后反引号内）
`gorm:"primaryKey"`      设置该字段为主键
`gorm:"column:xxx"`      指定数据库列名，不使用默认的字段名转蛇形
`gorm:"size:255"`        设置字符串长度
`gorm:"type:varchar(100)"` 指定数据库字段类型
`gorm:"not null"`        非空约束
`gorm:"unique"`          唯一约束
`gorm:"uniqueIndex"`     创建唯一索引
`gorm:"index:xxx"`       创建普通索引，可以给索引起名字
`gorm:"AUTO_INCREMENT"`  设置自增
`gorm:"-"`               完全忽略本字段，不映射数据库，仅Go内部使用

三、表名规则
1. 默认：结构体名 → 小写复数。例：People → peoples
2. 自定义方式1：给结构体实现 TableName() string 方法，GORM内部自动调用获取表名
3. 自定义方式2：db.Table("自定义表名").xxx()，代码中临时指定表名，优先级最高

四、数据库列名规则
默认：Go大驼峰字段名转为数据库蛇形小写。
例：UserName → user_name
如果写了 `gorm:"column:ani_name"`，则以column指定名字为准。

五、可空字段两种写法
1. 指针类型：Birthday *time.Time，nil代表数据库NULL
2. sql.NullXXX系列：Age sql.NullInt64，专门处理数据库可为NULL的基础类型

⚠️注意：AutoMigrate只会新增、修改字段，**不会删除数据库中已经存在的列**，生产环境不要直接依靠AutoMigrate做表变更。
*/

// 创建时间：2026/8/12 下午6:55
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Child struct {
	gorm.Model
	Name string `gorm:"default:'百家姓'"`    //默认值为空字符串，也可以指定默认值为百家姓
	Age  int    `gorm:"column:Child_age"` //在数据库中以man_age存在。默认值为0
}

// 内连接示例，仅演示语法，需要有parents表才能运行
type Parent struct {
	gorm.Model
	ChildID uint
	Phone   string
}
type ChildJoinResult struct {
	Child
	Phone string
}

// main函数外部，包级别
// Scopes示例：封装通用查询条件
func AgeGt11(db *gorm.DB) *gorm.DB {
	return db.Where("Child_age > ?", 11)
}
func OrderByAgeDesc(db *gorm.DB) *gorm.DB {
	return db.Order("Child_age desc")
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
	err = db.AutoMigrate(&Child{})
	/*
		疑问点：是不是db和Child绑定，代码就不能处理另一个表？
		答案：db对象本身不绑定任何表！
		前提：err = db.AutoMigrate(&Child{}, &Man{})，一次性把多张表完成建表。
		真正决定操作哪一张表的，是传给Find/First/Create的【结构体指针】：
			传入 &Child / []Child → 操作 children 表
			传入 &Man / []Man     → 操作 mans 表
		db只是通用数据库连接句柄，可以反复复用，同一个main函数来回切换不同模型做增删改查。
		注意：Model()方法会提前锁定模型，此时Find接收的结构体最好和Model模型保持一致。
	*/
	if err != nil {
		log.Fatal("自动迁移建表失败：", err)
	}
	fmt.Println("✅People表迁移完成")

	//创建记录：数据库中新增一行数据=结构体新增一个实例。在学习查询时，先创建一些数据，再注释掉。
	//child1 := Child{Name: "赵云", Age: 11,} //指定字段，则使用指定值
	//db.Create(&child1)
	//child2 := Child{Name: "张飞", Age: 12,}
	//db.Create(&child2)
	//child3 := Child{Name: "关羽", Age: 13,}
	//db.Create(&child3)

	//查询：一般查询
	//fmt.Println("==================================查询：一般查询====================================================")
	//var child1 Child //先创建一个结构体实例，用于接收查询结果
	//db.First(&child1) //根据主键，查询第一条记录，并保存到child1中。这里传的是结构体指针
	//fmt.Println("✅child1查询数据完成，child1=", child1)
	//var child2 Child
	//db.Last(&child2) //根据主键，查询最后一条记录，并保存到child2中
	//fmt.Println("✅child2查询数据完成，child2=", child2)
	//var child3 Child
	//db.Take(&child3) //随机查询一条记录，并保存到child3中
	//fmt.Println("✅child3查询数据完成，child3=", child3)
	//var child4 []Child //先创建一个结构体切片，用于接收多个查询结果
	//db.Find(&child4) //查询所有记录，并保存到child4中
	//fmt.Println("✅child4查询数据完成，child4=", child4)
	//var child5 Child
	//db.First(&child5,2) //查询id为2的记录，并保存到child5中。只有主键为整型时使用。等价于 select * from children where id=2;
	//fmt.Println("✅child5查询数据完成，child5=", child5)

	//查询：含where的查询
	//fmt.Println("==================================查询：含where的查询====================================================")
	//var child6 Child
	//db.Where("name=?","张飞").First(&child6) // key:value,得到第一个匹配的结果==> 查询name为张飞的记录，并保存到child6中。等价于 SELECT * FROM `children` WHERE name='张飞' AND `children`.`deleted_at` IS NULL ORDER BY `children`.`id` LIMIT 1;
	//fmt.Println("✅child6查询数据完成，child6=", child6)
	//var child7 []Child
	//db.Where("name=?","张飞").Find(&child7) // key:value,得到所有匹配的结果==> 查询name为张飞的记录，并保存到child7中。等价于 SELECT * FROM `children` WHERE name='张飞' AND `children`.`deleted_at` IS NULL;
	//fmt.Println("✅child7查询数据完成，child7=", child7)
	//var child8 []Child
	//db.Where("name<>?","赵云").Find(&child8) //key:value,得到所有匹配的结果==> 查询name不是张飞的所有记录，并保存到child8中。
	//fmt.Println("✅child8查询数据完成，child8=", child8)
	//var child9 []Child
	//db.Where("name in (?)",[]string{"张飞","赵云"}).Find(&child9) //key:value,得到所有匹配的结果==> 查询name是张飞或赵云的所有记录，并保存到child9中。
	//fmt.Println("✅child9查询数据完成，child9=", child9)
	//var child10 []Child
	//db.Where("name not in (?)",[]string{"张飞","赵云"}).Find(&child10) //key:value,得到所有匹配的结果==> 查询name不是张飞和赵云的所有记录，并保存到child10中。
	//fmt.Println("✅child10查询数据完成，child10=", child10)
	//var child11 []Child
	//db.Where("name like ?","%飞%").Find(&child11) //key:value,得到所有匹配的结果==> 查询name中包含飞的所有记录，并保存到child11中。
	//fmt.Println("✅child11查询数据完成，child11=", child11)
	//var child12 []Child
	//db.Where("name=? and age>12", "张飞",12).Find(&child12) //key:value,得到所有匹配的结果==> 查询name是张飞且age大于12的所有记录，并保存到child12中。
	//fmt.Println("✅child12查询数据完成，child12=", child12)
	//var child13 []Child
	//now := time.Now()
	//lastWeek := time.Now().AddDate(0,0,-7) // 当前时间往前推7天（上周）
	//db.Where("updated_at > ?",lastWeek).Find(&child13) //key:value,得到所有匹配的结果==> 查询updated_at大于lastWeek的所有记录，并保存到child13中。
	//fmt.Println("✅child13查询数据完成，child13=", child13)
	//var child14 []Child
	//db.Where("created_at between ? and ?",lastWeek,now).Find(&child14) //key:value,得到所有匹配的结果==> 查询created_at在lastWeek和now之间的所有记录，并保存到child14中。")
	//fmt.Println("✅child14查询数据完成，child14=", child14)

	//查询：Struct & Map查询
	//fmt.Println("==================================查询：Struct & Map查询====================================================")
	/*
		Where(结构体{})
		   - 仅使用【非零字段】作为查询条件；Go零值("",0,false,nil)直接忽略，不会生成SQL条件
		   - 坑：无法查询 age=0 / name="" 这类条件，零值会被丢弃
		   - 适合：全部条件都不为零的场景
	*/
	//var child15 Child
	//db.Where(&Child{Name:"张飞",Age:12}).First(&child15)
	//fmt.Println("✅child15查询数据完成，child15=", child15)
	/*
		Where(map[])
		map会保留零值，可以查询 age:0；适合动态组装多条条件
	*/
	//var child16 []Child
	//db.Where(map[string]interface{}{"name":"张飞","age":0}).Find(&child16)
	//fmt.Println("✅child16查询数据完成，child16=", child16)
	//主键的切片
	//var child17 []Child
	//db.Where([]int{3,2}).Find(&child17)//等价于sleect * from children where id in (3,2)

	//查询：or条件
	//fmt.Println("==================================查询：or条件====================================================")
	//var child18 []Child
	//db.Where("name=?","张飞").Or("name=?","赵云").Find(&child18) //等价于select * from children where name='张飞' or name='赵云' 会把名字是张飞，或者名字是赵云的全部行，全部查出来。
	//fmt.Println("✅child18查询数据完成，child18=", child18)

	//查询：内联条件
	/*
	   内联条件：不单独调用Where()，直接把查询条件写在【立即执行方法】(First/Find/Update/Delete等)的参数里。

	   特点：
	   1. 条件直接写在First、Find等函数参数中，不再链式调用Where()
	   2. ⚠️重要：内联条件**只作用于当前这一次立即执行方法**，不会传递、不会继承给后续链式调用的其他方法。
	   3. 支持多种格式：直接传主键值、sql片段+占位符、结构体对象。
	*/
	//fmt.Println("==================================查询：内联条件====================================================")
	//var child19 Child
	// 内联条件：直接传入数字，GORM默认把该值当做主键ID查询
	// 等价SQL：SELECT * FROM `children` WHERE `children`.`id` = 3 ORDER BY `children`.`id` LIMIT 1
	//db.First(&child19,3)
	//fmt.Println("✅child19查询数据完成，child19=", child19)
	//var child20 Child
	// 内联条件：SQL片段 + ?占位符参数
	// 适用于非数字主键、普通字段条件，防止SQL注入
	//db.First(&child20,"id=?","string_primary_key") //找 id 字段等于字符串 string_primary_key 的记录，最多返回 1 行。
	//fmt.Println("✅child20查询数据完成，child20=", child20)
	//var child21 Child
	// ⚠️注意坑：Find接收单个结构体变量，只会把第一条匹配记录赋值给child21
	// 内联条件：SQL片段占位符形式写条件
	//db.Find(&child21,"name=?","张飞")
	//fmt.Println("✅child21查询数据完成，child21=", child21)
	//var child22 []Child
	// 内联条件：传入结构体对象作为查询条件
	// 和Where(结构体)规则一样：**只会使用结构体的非零值作为查询条件，零值字段直接忽略**
	// 这里只取 Name:"赵云"，其他字段如果是零值不会生成where条件
	//db.Find(&child22,Child{Name:"赵云"})
	//fmt.Println("✅child22查询数据完成，child22=", child22)
	/*
		   补充总结：
		   1. 内联条件写法：条件直接写在First/Find的第二个参数位置，省略Where()链式调用
		   2. 支持三种传入格式
		      ① 直接传数值：db.First(&obj, 123) → 按主键ID查询
		      ② sql片段+占位符：db.First(&obj,"name=?","xxx")
		      ③ 结构体对象：db.Find(&slice, Child{Name:"赵云"})，仅非零字段生效
		   3. 关键区别链式Where：
		      Where链式的条件会向后传递给后续方法；
		      内联条件仅属于当前这个Find/First，不会向后传递。
		   4. 注意：Find传入结构体（非切片），只会拿到第一条匹配结果；多条数据必须用切片[]Child接收。
		   5. 结构体做内联条件同样存在零值陷阱：int=0、string=""这类零值不会生成查询条件。

		   额外举一个 “内联条件不会向后传递” 的示例
			内联条件只给First，后面的Count不受这个条件影响
			db.First(&child19, 3).Count(&total)
			Count()统计的是全表总数！不会带上 id=3 的条件，这就是内联条件不向后传递
	*/

	//FirstOrInit:获取匹配的第一条记录，否则根据给定的条件初始化一个新的对象（仅支持struct和map条件）
	//fmt.Println("==================================查询：FirstOrInit====================================================")
	//var child23 Child
	//db.FirstOrInit(&child23,Child{Name:"张飞",Age:12}) //查询name为张飞且age为12的第一条记录，如果没有找到，则初始化一个新的Child对象，并把name设置为张飞，age设置为12
	//fmt.Println("✅child23查询数据完成，child23=", child23)

	// Attrs：如果记录未找到，将使用Attrs里面的参数初始化struct；如果查到记录，Attrs内容直接被忽略
	//fmt.Println("==================================查询：Attrs====================================================")
	//var child24 Child
	// 逻辑：先用Where条件去数据库查询 Name="张飞" 的记录
	// ① 如果查询【找到了】记录 → Attrs(Child{Age:12}) 完全不生效，child24直接赋值数据库查出来的值
	// ② 如果查询【没找到】记录 → 不会插入数据库！仅仅在内存中把child24初始化为 Where条件 + Attrs合并后的结构体：Name="张飞", Age=12
	//db.Where(Child{Name:"张飞"}).Attrs(Child{Age:12}).First(&child24)
	//fmt.Println("✅child24查询数据完成，child24=", child24)
	/*
	   Attrs重点知识点：
	   1. Attrs 只会修改【内存中的结构体变量】，**不会向数据库新增任何数据！**
	      找不到记录也不会执行INSERT，仅仅把Go内存里的结构体填充上Attrs的字段。
	   2. 执行流程拆解：
	      第一步：Where(Child{Name:"张飞"}) 设置查询条件 name = "张飞"
	      第二步：First()执行SQL查询数据库
	         情况1：查到数据 → child24 = 数据库返回行，Attrs直接丢弃不使用
	         情况2：查不到数据(gorm.ErrRecordNotFound) → 内存child24 = Where条件合并Attrs的结构体{Name:"张飞",Age:12}，数据库依旧是空
	   3. 对比区别 Attrs vs Assign
	      - Attrs：仅在【查询不到】的时候才赋值内存结构体；查到则完全不用
	      - Assign：无论查询是否成功，都会把Assign里面字段覆盖到内存结构体上
	   4. 衍生方法：FirstOrCreate()
	      Attrs + FirstOrCreate：查不到的时候，不仅内存赋值，还会把这条记录真正插入数据库。
	      db.Where(Child{Name:"张飞"}).Attrs(Child{Age:12}).FirstOrCreate(&child24)
	   5. 注意Where传结构体同样有零值陷阱，Where里面结构体零值字段不会生成SQL条件。
	*/

	//高级查询：子查询
	/*
	   子查询：把一条查询嵌套在另一条SQL内部；
	   GORM中直接将 *gorm.DB 对象作为Where的参数，自动生成带括号的子查询SQL；
	   子查询只组装SQL，不会分开多次访问数据库，是一条SQL完成查询。
	   示例需求：查询年龄大于平均年龄的Child记录
	*/
	//fmt.Println("==================================查询：高级查询-子查询====================================================")
	//var child25 []Child
	// 先组装内层子查询：计算所有child的平均年龄，只是组装SQL，不执行
	/*
		db.Model(&Child{})：指定操作模型Child，对应表children，自动带上软删除 deleted_at IS NULL
		.Select("AVG(Child_age)")：指定要查询的字段，这里是聚合函数，计算 Child_age 的平均值
	*/
	//subQuery := db.Model(&Child{}).Select("AVG(Child_age)")
	// 主查询，把subQuery作为条件参数，GORM自动加()
	//db.Where("Child_age > (?)", subQuery).Find(&child25)
	//等价SQL：SELECT * FROM `children` WHERE `children`.`deleted_at` IS NULL AND Child_age > (SELECT AVG(Child_age) FROM `children` WHERE `children`.`deleted_at` IS NULL);
	//fmt.Println("✅child25子查询结果，child25=", child25)
	//var child26 []Child
	// 子查询搭配IN：查询id在子查询结果集中的数据
	/*
		db.Model(&Child{})：绑定 children 表，自动带上软删除条件
		.Select("id")：子查询只查询 id 这一列（IN 子查询一般只需要返回一列）
		.Where("Child_age > ?",11)：子查询自己的过滤条件，年龄大于 11
	*/
	//subQuery2 := db.Model(&Child{}).Select("id").Where("Child_age > ?",11)
	//db.Where("id IN (?)", subQuery2).Find(&child26)
	//等价SQL：SELECT * FROM `children` WHERE `children`.`deleted_at` IS NULL AND id IN (SELECT id FROM `children` WHERE `children`.`deleted_at` IS NULL AND Child_age > 11);
	//fmt.Println("✅child26子查询IN结果，child26=", child26)
	/*
	   子查询注意点：
	   1. subQuery只是组装查询链，没有调用Find/First，不会执行SQL；
	   2. 将subQuery传给Where，GORM自动处理括号，不要自己手动写括号；
	   3. 子查询和主查询在同一条SQL内完成，减少数据库交互；
	   4. 子查询可以用于 > < = IN 等场景，适合复杂统计过滤。
	*/

	//高级查询：指定想从数据库中检索出的字段，默认是查询全部字段
	//fmt.Println("==================================查询：高级查询-指定想从数据库中检索出的字段====================================================")
	//var child27 []Child
	//db.Select("id","name").Find(&child27) //等价SQL：SELECT id,name FROM `children` WHERE `children`.`deleted_at` IS NULL;
	//fmt.Println("✅child27指定想从数据库中检索出的字段，child27=", child27)
	//var child28 []Child
	//db.Select("name").Find(&child28) //等价SQL：SELECT name FROM `children` WHERE `children`.`deleted_at` IS NULL;

	//高级查询：排序
	/*
	   Order：指定数据库返回记录的排序规则
	   链式多次调用Order：默认是追加排序条件；
	   Order(条件, true)：第二个参数传true，代表覆盖，丢弃前面所有已经设置的排序，只用当前这一条排序。
	   语法：Order("列名 排序方式")  desc降序，asc升序
	*/
	//fmt.Println("==================================查询：高级查询‑排序====================================================")
	//var child29 []Child
	// 一次Order写多个字段：先按name升序，name相同再按Child_age降序
	//等价SQL：SELECT * FROM `children` WHERE `children`.`deleted_at` IS NULL ORDER BY name,Child_age desc;
	//db.Order("name ,Child_age desc").Find(&child29)
	//fmt.Println("✅child29排序，child29=", child29)
	//var child30 []Child
	// 多次链式Order，默认【追加】：先Child_age desc，再name升序
	//等价SQL：SELECT * FROM `children` WHERE `children`.`deleted_at` IS NULL ORDER BY Child_age desc,name;
	//db.Order("Child_age desc").Order("name").Find(&child30)
	//fmt.Println("✅child30排序，child30=", child30)
	/*
		//下面两句完全等价！多次Order就是追加拼接
		db.Order("a desc, b asc").Find(&x)
		db.Order("a desc").Order("b asc").Find(&x)
		两种Order写法规则：
		1. 单Order内逗号分隔、多次链式Order追加，规则完全相同；
		2. 书写的先后顺序 = MySQL排序优先级；写在前的优先排序；
		3. 只要【字段顺序+升降序】保持一致，两种写法等价；顺序不一样，则排序逻辑不一样。
	*/
	//var child31 []Child
	//var child32 []Child
	// 注意：gorm链式对象会复用！
	// 第一步：设置排序 Child_age desc，执行Find查询得到child31
	//新gorm以移除：db.Order("Child_age desc").Find(&child31).Order("name", true).Find(&child32)
	// Order第二个参数true → 覆盖模式！丢弃之前 Child_age desc，只保留新条件 name
	/*
	   排序知识点总结：
	   1. Order("col desc, col2 asc")：一个Order内写多个排序字段，逗号分隔。
	   2. 多次调用Order(条件)，不传第二个参数：条件向后追加。
	   3. Order(条件, true)：true代表覆盖，清空之前全部排序规则，仅使用当前条件。  新gorm以移除
	   ⚠️坑：链式*gorm.DB是可变对象，调用Find后，链式对象还可以继续复用，修改条件继续查询。
	   child31、child32是两次独立数据库查询，不是内存切片处理排序，排序逻辑发生在MySQL层。
	*/

	//高级查询：指定从数据库检索出的最大记录数
	//fmt.Println("==================================查询：高级查询‑指定从数据库检索出的最大记录数====================================================")
	//var child33 []Child
	//db.Limit(3).Find(&child33) //等价SQL：SELECT * FROM `children` WHERE `children`.`deleted_at` IS NULL LIMIT 3;
	//fmt.Println("✅child33查询结果，child33=", child33)
	//var child34 []Child
	//db.Limit(-1).Find(&child34) //填-1，表示不限制，查询全部
	//fmt.Println("✅child34查询结果，child34=", child34)

	//高级查询：Offset / Limit 分页
	/*
	   Offset(n)：跳过前面 n 条记录，从第 n+1 条开始取数据；常用于分页。
	   Limit(n)：最多只返回 n 条记录。
	   特殊：Limit(-1) 代表清除之前设置的Limit限制，返回全部数据。
	   注意：Offset单独使用没有数量限制，会返回Offset之后所有符合条件的行；生产务必搭配Limit防止返回海量数据。
	*/
	//fmt.Println("==================================查询：高级查询‑Offset Limit分页====================================================")
	//var child35 []Child
	// Offset(3)：跳过前3条记录，从第4条开始拿，没有Limit会拿后面全部数据
	//等价SQL：SELECT * FROM `children` WHERE `children`.`deleted_at` IS NULL OFFSET 3;
	//db.Offset(3).Find(&child35)
	//fmt.Println("✅child35查询结果，跳过前3条，取后续全部：", child35)
	//var child36 []Child
	// Limit(-1)：清除链式上已经设置的Limit限制，返回全部匹配记录
	// 场景：前面链式设置过Limit，后续想取消限制，就用Limit(-1)
	//db.Limit(-1).Find(&child36)
	//fmt.Println("✅child36查询结果，清除limit限制，返回全部：", child36)
	/*
		标准分页公式：Offset((page‑1)*pageSize).Limit(pageSize).Find(&list)
		page：页码，从1开始
		pageSize：一页多少条
		Offset(n)：丢弃前面n行；Limit(n)：最多读取n行
		示例 page=2，pageSize=2：
		Offset((2‑1)*2) → Offset(2)，跳过前2条；Limit(2)读取2条，拿到第2页数据。
		⚠️注意：
		1. MySQL的Limit offset分页，数据量大的时候性能会变差；
		2. 分页一定要同时写Offset+Limit，不要只写Offset。
	*/

	//高级查询：Count 统计记录总数
	/*
	   Count(&变量)：统计符合条件的行数，结果存入传入的int变量
	   注意坑：
	   1. Count参数必须传 **int64指针**，不要用int；
	   2. 如果写在Find()之后链式调用，Find已经把条件带入，Count会继承Where条件；
	   3. Table("表名")可以直接指定表，不依赖Model模型；
	   4. Select搭配Count：Count()会覆盖Select设置的字段，内部强制使用COUNT(*)。
	*/
	//fmt.Println("==================================查询：高级查询‑Count====================================================")
	//var total int64 // count统计结果接收变量，必须int类型
	//var child37 []Child
	// 先Find查询出age=11的数据存入child37；再链式Count，继承Where条件，统计age=11的总条数
	//等价SQL：SELECT * FROM children WHERE deleted_at IS NULL AND Child_age = 11;
	//         SELECT count(*) FROM children WHERE deleted_at IS NULL AND Child_age = 11;
	//db.Where("Child_age = ?", 11).Find(&child37).Count(&total)
	//fmt.Println("✅child37查询结果，child37=", child37)
	//fmt.Println("✅满足age=11的记录总数 total =", total)

	// Table直接指定表名，统计整张表有效总条数（自动带上软删除 deleted_at IS NULL）
	//等价SQL：SELECT count(*) FROM `children` WHERE `children`.`deleted_at` IS NULL;
	//db.Table("children").Count(&total)
	//fmt.Println("✅children表全部有效记录总数 total =", total)

	// ⚠️坑：写Select("count(name)")，调用Count后会被覆盖，实际执行依旧是 count(*)
	//等价SQL：SELECT count(*) FROM `children` WHERE `children`.`deleted_at` IS NULL;
	//db.Table("children").Select("count(name)").Count(&total)
	//fmt.Println("✅即使写Select(count(name))，Count依旧执行count(*), total =", total)
	/*
	   Count知识点总结：
	   1. Count会继承链式上的Where条件；写在Find之后，会额外再发一条统计SQL，做两次数据库请求。
	   2. Count() 内部固定使用 COUNT(*)，会把前面的Select覆盖掉。想要自定义count(列)，要用Row/Raw原生SQL。
	   3. 嵌入gorm.Model模型，Count自动过滤软删除deleted_at不为NULL的数据。
	*/
	// 拓展：想要执行 count(name) 自定义统计，不能用Count()方法，用原生Raw
	//db.Raw("SELECT count(name) FROM children WHERE deleted_at IS NULL").Scan(&total)

	//高级查询：Group 分组和having  比较复杂

	//高级查询：Join 连接查询
	/*
	   Join：多表关联查询
	   - Join：内连接，只取两边匹配的数据
	   - LeftJoin：左连接，左表全部保留，右表匹配不上为null
	   语法：Join("关联表 别名","关联条件")
	   注意：关联查询接收结果建议使用自定义结构体，直接用Model结构体容易字段映射错乱。
	*/
	//fmt.Println("==================================查询：高级查询‑Join连接查询====================================================")
	//var joinList []ChildJoinResult
	// ✅内连接，使用Joins，字符串写完整 inner join
	//db.Model(&Child{}).
	//	Joins("INNER JOIN parents p ON p.child_id = children.id").
	//	Find(&joinList)
	//fmt.Println("✅内连接查询结果：", joinList)

	//var leftJoinList []ChildJoinResult
	// ✅左连接，使用Joins，字符串写完整 left join
	//db.Model(&Child{}).
	//	Joins("LEFT JOIN parents p ON p.child_id = children.id").
	//	Find(&leftJoinList)
	//fmt.Println("✅左连接查询结果：", leftJoinList)
	/*
	   Join小知识点：
	   1. Join/LeftJoin第一个参数是表名+别名；第二个参数是on关联条件；
	   2. 多表查询不要直接用原始Model接收，建议定义联合结果结构体；
	   3. 内连接Join：两边条件匹配才返回；左连接LeftJoin：左表全部保留，右表无匹配字段为NULL；
	   4. AutoMigrate需要把被关联模型也执行迁移 db.AutoMigrate(&Child{},&Parent{})
	*/

	//高级查询：Pluck 查询model中的一个列作为切片，如果想要查询多个列，应该使用Scan
	/*
	   Pluck(列名, &切片变量)
	   作用：执行SQL，**只查指定的这一列**，把该列全部行的值直接填充到切片里。
	   不需要完整结构体接收，直接拿到一列数据。
	   关键点：
	   1. Pluck 本身就会执行数据库查询，不需要提前调用Find。
	   2. 仅支持**单个列**，不能一次性拿多列；多列要用Scan。
	   3. 接收的切片类型必须和数据库字段类型匹配：数字用[]int64，字符串用[]string。
	   4. 链式继承前面的条件：Where / Model / Limit / Offset 都会生效。
	   5. 如果前面写了Find再调用Pluck：Find执行一次查询，Pluck会**再执行一次全新SQL**，两次查询互不干涉，会多一次数据库访问。
	*/
	//fmt.Println("==================================查询：高级查询‑Pluck查询model中的一个列作为切片====================================================")
	//var child38 []Child
	//var ages []int64
	// Find：执行一次SQL把全部记录读到child38；接着Pluck又执行一条全新SQL，取出Child_age列放入ages
	// ⚠️两次独立SQL，child38的数据和ages互不依赖，此处会产生多余数据库IO，不推荐该写法
	//db.Find(&child38).Pluck("Child_age", &ages)
	//fmt.Println("✅child38查询结果，ages=", ages)
	//var name []string
	//db.Model(&Child{}).Pluck("name", &name)
	//fmt.Println("✅child38查询结果，name=", name)
	//db.Table("children").Pluck("name", &name)
	//var child39 []Child
	// Select指定查询多列，Find存入Child切片；没有被select查询出来的结构体字段，保持Go零值
	//db.Select("child_age,name").Find(&child39)
	//fmt.Println("✅child39查询结果，child39=", child39)

	//高级查询：Scopes 查询范围,建立在链式操作的基础上。基于此，可以抽取一些通用逻辑，包装成可重用的函数库
	/*
	   Scopes 作用：把通用的查询逻辑（Where、Order、Limit等链式条件）封装成函数，实现查询逻辑复用。
	   1. 封装函数签名必须是：func(db *gorm.DB) *gorm.DB，输入输出都是*gorm.DB
	   2. 在db.Scopes(函数名) 传入封装好的函数，会把里面的查询条件合并到当前链式
	   3. 可以同时传入多个scope函数，条件会依次叠加
	   4. 适合：通用过滤、公共排序、通用分页逻辑，多处复用，减少重复代码
	*/
	//fmt.Println("==================================查询：高级查询‑Scopes查询范围====================================================")
	//var child40 []Child
	// 使用Scopes，复用封装好的查询条件
	//db.Model(&Child{}).
	//	Scopes(AgeGt11, OrderByAgeDesc).Find(&child40) // 传入多个scope，条件依次叠加。
	/*
		AgeGt11的作用：组装查询条件，Where("Child_age > ?", 11)
		OrderByAgeDesc的作用：组装查询条件，Order("Child_age desc")
	*/
	//fmt.Println("✅Scopes复用通用条件查询结果 child40 =", child40)
	/*
	   Scopes注意点：
	   1. scope函数不能写在main函数内部，需要定义在包级别；
	   2. scope只组装查询条件，**不会执行SQL**；真正执行还是靠Find/First等立即执行方法；
	   3. 多个Scopes函数之间条件是叠加关系；
	   4. 可以结合变量，做动态scope，实现通用分页等逻辑。
	*/
	// 拓展：带参数的Scope（闭包写法，写在main内部）
	/*
	   ageFilter := func(minAge int) func(db *gorm.DB) *gorm.DB {
	   	return func(db *gorm.DB) *gorm.DB {
	   		return db.Where("Child_age > ?", minAge)
	   	}
	   }
	   var child41 []Child
	   db.Model(&Child{}).
	   	Scopes(ageFilter(12)). // 传入参数，动态生成查询条件
	   	Find(&child41)
	   fmt.Println("✅带参数闭包scope结果 child41 =", child41)
	*/

	//高级查询：多个立即执行函数:后一个立即执行方法Count会复用前一个立即执行方法Find的条件，不包括内联条件（"Child_age > ?", 11）
	var count int64
	var childs []Child
	db.Where("name = ?", "赵云").Find(&childs, "Child_age > ?", 11).Count(&count) //相当于sql的 select * from children where name = '赵云'

}

/*
立即执行方法：指能立即生成Sql语句并发送到数据库的方法，一般是crud方法，比如：Create、Update、Delete、Find、First、FirstOrCreate、Pluck、Count、Row、Scan、Table、Raw
*/

learn_layer_demo 项目使用与调用链路手册
一、项目简介
本项目是基于 Go+Gin+GORM 的三层架构练习 Demo，无前端页面，仅提供 HTTP 接口，核心目的是熟练掌握 Handler-Service-Repository 分层规范与单向调用流程。
技术栈：Go 语言 + Gin Web 框架 + GORM ORM + MySQL
核心功能：用户模块完整 CRUD（增、删、改、查）
定位：分层架构入门练习项目，也可作为后续业务项目的基础脚手架
设计原则：严格遵循「Handler 管请求响应、Service 管业务规则、Repository 管数据库操作」的单向调用规范
二、目录与文件说明
plaintext
learn_layer_demo/
├── go.mod / go.sum    # Go 模块依赖管理文件
├── main.go            # 程序总入口：初始化组件 + 注册路由 + 启动服务
├── config/
│   └── db.go          # MySQL 连接初始化 + GORM 自动迁移建表
├── model/
│   └── user.go        # 用户数据模型，对应数据库 users 表结构
├── handler/
│   └── user_handler.go # 接口控制层：接收请求、参数解析、返回响应
├── service/
│   └── user_service.go # 业务逻辑层：业务校验、流程编排
├── repository/
│   └── user_repo.go    # 数据访问层：封装数据库增删改查
├── middleware/         # 预留目录，后续存放鉴权、跨域等中间件
└── utils/              # 预留目录，后续存放通用工具函数
三、启动运行完整步骤
前置准备
本地已安装 Go 开发环境、MySQL 服务
提前在 MySQL 中手动创建数据库（GORM 不会自动建库）：
sql
CREATE DATABASE learn_layer_demo DEFAULT CHARACTER SET utf8mb4;
启动流程
修改数据库配置
打开 config/db.go，将 dsn 中的账号、密码、数据库名修改为你本地 MySQL 的实际配置：
go
运行
dsn := "root:你的密码@tcp(127.0.0.1:3306)/learn_layer_demo?charset=utf8mb4&parseTime=True&loc=Local"
安装依赖
在项目根目录打开终端，执行：
bash
go mod tidy
启动项目
bash
go run main.go
验证启动成功
终端输出 Listening and serving HTTP on :8080 即为启动成功。程序启动时会自动执行 AutoMigrate，在数据库中创建 users 表。
四、核心接口完整调用链路
所有接口严格遵循单向调用：客户端 → Handler → Service → Repository → MySQL → 原路返回
1. 根据 ID 查询用户（GET /user/:id）
   客户端发送 GET 请求，地址示例：http://127.0.0.1:8080/user/1
   Gin 匹配路由规则，将本次请求的全部信息封装到 *gin.Context，触发执行 handler.GetUserById
   Handler 层执行逻辑
   定义 UriParam 结构体，通过 c.ShouldBindUri(&param) 将路径中的 id 参数绑定到结构体
   参数绑定失败 → 直接返回 400 状态码 +「参数错误」提示
   绑定成功 → 提取 id 值，调用 service.GetUserInfo(id) 进入业务层
   Service 层执行逻辑
   执行业务校验：id 不能等于 0
   校验不通过 → 返回空的 model.User 结构体 + 业务错误信息
   校验通过 → 调用 repository.GetUserById(id) 进入数据层
   Repository 层执行逻辑
   声明 model.User 类型变量，用于接收查询结果
   执行 GORM 的 First 方法：按主键 id 查询 users 表第一条匹配记录
   查询数据填充到变量，同时返回数据库执行错误
   将「用户数据 + 错误」一同返回给 Service 层
   结果沿原路逐层向上返回，最终由 Handler 组装成 JSON 格式返回给客户端
2. 新增用户（POST /user）
   客户端发送 POST 请求，请求体携带 JSON 格式的姓名、年龄
   Gin 匹配路由，触发执行 handler.CreateUser
   Handler 层
   通过 c.ShouldBindJSON(&user) 将请求体解析为 User 结构体
   解析失败 → 返回 400 参数解析失败
   解析成功 → 调用 service.AddUser(&user)
   Service 层
   业务校验：年龄不能为负数
   校验失败 → 返回对应错误
   校验通过 → 调用 repository.CreateUser(user)
   Repository 层
   执行 GORM 的 Create 方法，向 users 表插入一条用户记录
   返回数据库写入错误
   结果原路返回，Handler 向客户端返回新增成功 / 失败的 JSON 响应
3. 更新用户（PUT /user）
   客户端发送 PUT 请求，请求体携带包含 ID、姓名、年龄的 JSON
   Gin 匹配路由，触发执行 handler.UpdateUser
   Handler 层
   解析请求体 JSON 到 User 结构体
   解析失败 → 返回 400 参数错误
   解析成功 → 调用 service.UpdateUser(&user)
   Service 层
   业务校验：用户 ID 不能为 0、年龄不能为负数
   校验失败 → 返回对应错误
   校验通过 → 调用 repository.UpdateUser(user)
   Repository 层
   执行 GORM 的 Save 方法，根据主键 ID 全量更新整条记录
   返回数据库更新错误
   结果原路返回，Handler 返回更新结果
4. 删除用户（DELETE /user/:id）
   客户端发送 DELETE 请求，地址示例：http://127.0.0.1:8080/user/5
   Gin 匹配路由，触发执行 handler.DeleteUser
   Handler 层
   解析路径参数 id
   解析失败 → 返回 400 参数错误
   解析成功 → 调用 service.DeleteUser(id)
   Service 层
   业务校验：id 不能为 0
   校验失败 → 返回错误
   校验通过 → 调用 repository.DeleteUser(id)
   Repository 层
   执行 GORM 的 Delete 方法；因结构体嵌入 gorm.Model 自带软删除，仅填充 deleted_at 字段，不会物理删除记录
   返回数据库执行错误
   结果原路返回，Handler 返回删除成功 / 失败提示
   五、接口测试参考表
   表格
   接口功能	请求方式	请求地址	请求示例 / 说明
   查询单个用户	GET	http://127.0.0.1:8080/user/1	路径末尾传用户 ID
   新增用户	POST	http://127.0.0.1:8080/user	请求体：{"name":"张三","age":20}
   更新用户	PUT	http://127.0.0.1:8080/user	请求体：{"ID":1,"name":"张三改","age":25}
   删除用户	DELETE	http://127.0.0.1:8080/user/5	路径末尾传用户 ID
   测试工具推荐使用 Postman、Apifox 等接口调试工具；浏览器地址栏仅能测试 GET 接口。
   六、注意事项
   软删除机制：User 结构体嵌入了 gorm.Model，删除操作为软删除，记录不会从数据库物理消失，普通查询无法查到已删除数据。
   自动建表规则：程序启动时仅自动创建不存在的数据表，不会覆盖已有数据；数据库本身需要手动提前创建。
   分层红线：禁止 Handler 直接操作数据库，禁止 Service 层接收 gin.Context，禁止 Repository 层写业务逻辑，始终保持单向调用。
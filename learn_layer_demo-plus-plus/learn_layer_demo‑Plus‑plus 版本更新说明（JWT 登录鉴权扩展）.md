# learn_layer_demo‑Plus‑plus 版本更新说明（JWT 登录鉴权扩展）

## 一、版本概述

本版本基于 Plus V3（DTO/VO + validator 参数校验 + 统一响应封装）继续迭代，在原有完整分层 CRUD 脚手架基础上，新增**JWT 身份认证能力**，实现登录下发令牌、接口鉴权拦截整套流程。
原有 Handler‑Service‑Repository 分层、DTO‑VO 数据模型、参数校验、统一响应全部保留，仅新增登录与鉴权相关模块，不改动原有用户增删改查业务逻辑。

## 二、本次版本新增解决的问题

Plus‑V3 版本所有接口全部开放访问，没有登录身份校验，任何人都可以直接调用全部接口。
本版本引入 JWT 鉴权机制：

1. 提供登录接口，账号密码校验通过后下发 JWT 令牌给前端；
2. 区分公开接口、受保护鉴权接口；受保护接口请求必须携带合法 Bearer Token；
3. 中间件统一拦截鉴权，token 缺失、篡改、过期直接返回 401 业务错误，不进入业务 handler；
4. 鉴权解析出登录用户 ID，存入 gin 上下文，业务接口从上下文获取当前登录用户标识。

>
> 原有 CRUD 接口逻辑、DTO/VO、validator 校验、utils 统一响应工具全部复用。

## 三、本次学习目标

1. 理解 JWT 令牌生成、下发、客户端携带、服务端解析鉴权完整链路；
2. 掌握 Gin 路由分组结合中间件：部分接口公开，部分接口全局鉴权；
3. 理解 Gin 上下文传递鉴权解析后的 userId 给下游 handler；
4. 分清登录公开接口、受保护接口的路由分组配置；
5. 在现有分层架构上扩展鉴权组件，不破坏原有分层职责边界。

## 四、新增依赖

```
go get github.com/golang-jwt/jwt/v5
```

## 五、新增与改动内容明细

### 1. 新增文件

1. `dto/login_dto.go`：新增登录请求 DTO 结构体，搭配 validate 校验标签，接收登录账号密码入参。
2. `middleware/jwt_auth.go`：JWT 鉴权中间件，读取 Authorization 请求头，解析、校验 token，将 userId 存入 gin 上下文，鉴权失败直接拦截请求。
3. `utils/jwt_util.go`：JWT 工具包，包含`GenerateToken`生成 token、`ParseToken`解析 token 函数，定义自定义声明结构体`MyClaims`，维护 JWT 密钥、过期时长常量。

### 2. 修改文件

1. `handler/user_handler.go`
   - 新增`Login`登录 handler；新增`GetProfile`获取当前登录用户信息 handler；
   - Login 内部调用 service 登录逻辑，拿到 userId，调用 utils 生成 token，通过统一响应返回 token 给前端；
   - GetProfile 从 gin 上下文取出鉴权中间件存放的 userId，调用 service 获取用户个人信息，返回 VO；
2. `service/user_service.go`
   - 新增`Login`业务函数：接收用户名密码，调用 repository 查询用户，内存比对密码，返回用户 ID；
   - 新增`GetProfile`业务函数：接收 userId，调用 repository 查询用户，model 转 VO，屏蔽密码敏感字段返回；
3. `repository/user_repo.go`
   - 新增`GetUserByUsername`：根据用户名查询用户记录；
   - 新增`GetUserByIdPtr`：根据用户 ID 获取用户指针模型，供 profile 业务使用；
4. `main.go`
   - 调整路由分组：`apiGroup`为基础`/api`前缀；
   - 登录接口注册在`apiGroup`下，公开无鉴权；
   - 新建`authGroup`子路由组，挂载`middleware.JWTAuthMiddleware`中间件；
   - 将`/api/profile`注册到鉴权分组内，访问该接口必须携带 JWT 令牌。

### 3. 无改动文件

- `vo/user_vo.go`、`model/user.go`：原有结构体不变；
- 原有用户增删改查 handler/service/repository 代码完全保留；
- `utils/validator.go`、`utils/response.go`：参数校验、统一响应工具无修改，鉴权失败直接复用`utils.Fail()`输出统一格式 JSON。

## 六、接口能力变化

1. 公开接口（无需携带 token）
   - `POST /api/login`：登录接口，传入用户名密码，成功返回 JWT token；
   - 原有`/api/user`相关增删改查接口保持原有访问逻辑（本 demo 未做鉴权保护）。
2. JWT 鉴权保护接口（请求头必须携带 `Authorization: Bearer <token>`）
   - `GET /api/profile`：获取当前登录用户个人资料，从 token 解析 userId 查询用户信息。

## 七、新旧版本对比

表格

| 对比项 | Plus‑V3 | Plus‑plus (JWT 扩展版本) |
| --- | --- | --- |
| 接口访问控制 | 全部接口直接开放访问，无身份校验 | 区分公开接口、鉴权保护接口；受保护接口需要携带 Bearer Token |
| 登录能力 | 无登录逻辑 | 新增登录接口，账号密码校验后下发 JWT 令牌 |
| 身份识别 | 无用户身份概念 | JWT 中间件解析 token，userId 存入 gin.Context，handler 读取使用 |
| 请求拦截 | 仅做参数格式校验、业务校验 | 增加鉴权拦截；token 非法 / 缺失 / 过期直接返回 401，不执行业务逻辑 |
| 分层架构 | Handler‑Service‑Repository + DTO‑VO + validator + 统一响应 | 完整继承 V3 全部能力，新增 middleware、jwt 工具，原有分层职责不变 |

## 八、当前版本局限

1. 密码采用明文比对，真实项目应替换 bcrypt 哈希密码校验；
2. JWT 密钥硬编码在代码中，生产环境应迁移至配置文件；
3. JWT 无刷新 token、黑名单退出登录逻辑；
4. 仅对`/api/profile`开启鉴权，其余 CRUD 接口未接入鉴权分组。
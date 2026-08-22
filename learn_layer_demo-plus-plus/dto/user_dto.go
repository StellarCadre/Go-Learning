// 创建时间：2026/8/18 下午6:12
package dto

// CreateUserReq 新增用户-请求参数
type CreateUserReq struct {
	Name string `json:"name" validate:"required,min=1,max=20"` // 不能为空，1‑20字符
	Age  int    `json:"age" validate:"gte=0,lte=150"`          // 0~150
}

// UpdateUserReq 更新用户-请求参数
type UpdateUserReq struct {
	ID   uint   `json:"ID" validate:"required,gt=0"` // ID必须大于0
	Name string `json:"name" validate:"required,min=1,max=20"`
	Age  int    `json:"age" validate:"gte=0,lte=150"`
}

// UserIdUri 路径ID参数
type UserIdUri struct {
	ID uint `uri:"id" validate:"required,gt=0"`
}

// LoginDTO 用户登录请求参数
type LoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

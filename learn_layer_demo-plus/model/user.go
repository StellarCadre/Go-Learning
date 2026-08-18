// 创建时间：2026/8/17 下午6:57
package model

import "gorm.io/gorm"

// User 用户数据表模型
type User struct {
	gorm.Model
	Name string `gorm:"size:50;not null"`
	Age  int
}

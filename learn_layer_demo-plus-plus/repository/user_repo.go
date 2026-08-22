// 创建时间：2026/8/17 下午6:57
package repository

import (
	"learn_layer_demo/config"
	"learn_layer_demo/model"
)

// GetUserById 查询单个用户，无业务逻辑
func GetUserById(id uint) (model.User, error) {
	var user model.User //声明一个结构体变量 user，类型是User。 此时它是一个空结构体，用来存放一会儿从数据库查出来的结果。
	err := config.DB.First(&user, id).Error
	/*
		config.DB：初始化好的全局数据库连接对象
		.First(&user, id)：去 users 表里，按主键 id 查询第一条记录。
	*/
	return user, err
}

// CreateUser 新增用户
func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

// UpdateUser 更新用户（全量更新）
func UpdateUser(user *model.User) error {
	return config.DB.Save(user).Error
}

// DeleteUser 根据ID删除用户（软删除）
func DeleteUser(id uint) error {
	return config.DB.Delete(&model.User{}, id).Error
}

// ==========新增：JWT登录需要，包级函数，统一使用全局 config.DB==========
// GetUserByUsername 根据用户名查询
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetUserByIdPtr 根据id查询，返回指针版本，供登录业务使用
func GetUserByIdPtr(id uint) (*model.User, error) {
	var user model.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

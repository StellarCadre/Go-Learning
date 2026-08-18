// 创建时间：2026/8/17 下午7:00
package service

import (
	"errors"
	"learn_layer_demo/model"
	"learn_layer_demo/repository"
)

// GetUserInfo 获取用户，增加业务校验
func GetUserInfo(id uint) (model.User, error) {
	/*
		model.User 就是你在 model/user.go 里定义的 User 结构体。
		因为它定义在 model 包下，其他包（service、handler）要使用它，必须写成 model.User（包名。结构体名）。
		第一个返回值：查询到的用户数据，类型是 model.User 结构体；
		第二个返回值：错误信息，类型是 error。
	*/

	// 业务规则：id不能等于0
	if id == 0 {
		return model.User{}, errors.New("用户id非法") //因为业务校验不通过，没有有效数据，所以返回「空结构体 + 错误信息」。
	}
	return repository.GetUserById(id) //直接把 repository 层返回的「用户数据 + 错误」原封不动向上返回给 handler。
}

// AddUser 新增用户业务处理
func AddUser(user *model.User) error {
	// 业务规则：年龄不能小于0
	if user.Age < 0 {
		return errors.New("年龄不能为负数")
	}
	return repository.CreateUser(user)
}

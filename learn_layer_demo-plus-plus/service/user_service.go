// 创建时间：2026/8/17 下午7:00
package service

import (
	"errors"
	"learn_layer_demo/dto"
	"learn_layer_demo/model"
	"learn_layer_demo/repository"
	"learn_layer_demo/vo"
)

// GetUserInfo 查询用户：入参id，返回VO
func GetUserInfo(id uint) (vo.UserInfoVO, error) {
	// 业务校验
	if id == 0 {
		return vo.UserInfoVO{}, errors.New("用户id非法")
	}
	// 调用数据层，拿到Model
	user, err := repository.GetUserById(id)
	if err != nil {
		return vo.UserInfoVO{}, err
	}
	// Model 转 VO
	userVO := vo.UserInfoVO{
		ID:        user.ID,
		Name:      user.Name,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
	}
	return userVO, nil
}

// AddUser 新增用户：入参是DTO
func AddUser(req *dto.CreateUserReq) error {
	// DTO 转 Model
	userModel := &model.User{
		Name: req.Name,
		Age:  req.Age,
	}
	// 调用数据层写入
	return repository.CreateUser(userModel)
}

// UpdateUser 更新用户：入参是DTO
func UpdateUser(req *dto.UpdateUserReq) error {
	// 1.先查询数据库中该用户原始记录
	oldUser, err := repository.GetUserById(req.ID)
	if err != nil {
		return errors.New("用户不存在")
	}
	// 2.只修改需要变更的字段，其余保留数据库原来的值
	oldUser.Name = req.Name
	oldUser.Age = req.Age
	// 3.传入查询出来的对象，Save，CreatedAt保留原有值，GORM自动刷新UpdatedAt
	return repository.UpdateUser(&oldUser)
}

// DeleteUser 删除用户：入参id
func DeleteUser(id uint) error {
	// 业务校验
	if id == 0 {
		return errors.New("用户id非法")
	}
	// 调用数据层删除
	return repository.DeleteUser(id)
}

//=================JWT登录相关=================

// Login 校验账号密码 返回用户ID
func Login(username, password string) (uint, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return 0, err
	}
	// 真实项目后续替换为bcrypt加密密码比对
	if user.Password != password {
		return 0, errors.New("密码错误")
	}
	return user.ID, nil
}

// GetProfile 根据id查询用户，并转为VO返回（不返回密码字段）
func GetProfile(userId uint) (*vo.UserInfoVO, error) {
	user, err := repository.GetUserByIdPtr(userId)
	if err != nil {
		return nil, err
	}
	userVO := &vo.UserInfoVO{
		ID:   user.ID,
		Name: user.Name,
		Age:  user.Age,
	}
	return userVO, nil
}

/*
对service.UserService.Login和service.UserService.GetProfile的分析：


*/

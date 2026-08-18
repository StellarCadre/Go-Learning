// 创建时间：2026/8/17 下午7:01
package handler

import (
	"github.com/gin-gonic/gin"
	"learn_layer_demo/dto"
	"learn_layer_demo/service"
	"learn_layer_demo/utils"
)

// GetUserById  查询用户接口
func GetUserById(c *gin.Context) {
	// 用DTO里的结构体接收路径参数。   ShouldBindUri：把url路径参数解析赋值给param
	var param dto.UserIdUri
	if err := c.ShouldBindUri(&param); err != nil {
		//旧：c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		//新： 使用utils封装的Fail方法
		utils.Fail(c, 400, 40000, "参数解析失败")
		return
	}

	//新： 执行validator参数格式校验
	if err := utils.Validator.Struct(&param); err != nil {
		//旧：c.JSON(http.StatusBadRequest, gin.H{"msg": "参数格式不合法"})
		//新：
		utils.FailWithValidator(c, err)
		return
	}
	// 调用Service，拿到VO
	userVO, err := service.GetUserInfo(param.ID)
	if err != nil {

		//旧：c.JSON(http.StatusNotFound, gin.H{"msg": "查询失败:" + err.Error()})
		utils.Fail(c, 404, 50000, "查询失败:"+err.Error())
		return
	}
	// 直接返回VO给前端
	//旧：c.JSON(http.StatusOK, gin.H{"data": userVO})
	utils.Success(c, userVO)
}

// CreateUser 新增用户接口
func CreateUser(c *gin.Context) {
	// 用DTO接收请求体，不再绑定Model    ShouldBindJSON：把http请求体json解析赋值给req
	var req dto.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		//旧：c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析失败"})
		utils.Fail(c, 400, 40000, "参数解析失败")
		return
	}

	//新： 执行validator参数格式校验
	if err := utils.Validator.Struct(&req); err != nil {
		//旧：c.JSON(http.StatusBadRequest, gin.H{"msg": "参数格式不合法"})
		utils.FailWithValidator(c, err)
		return
	}
	err := service.AddUser(&req)
	if err != nil {
		//旧：c.JSON(http.StatusInternalServerError, gin.H{"msg": "新增失败:" + err.Error()})
		utils.Fail(c, 500, 50000, "新增失败:"+err.Error())
		return
	}
	//旧：c.JSON(http.StatusOK, gin.H{"msg": "新增成功"})
	utils.Success(c, nil)
}

// UpdateUser 更新用户接口
func UpdateUser(c *gin.Context) {
	// 用DTO接收请求体
	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		//旧：c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析失败"})
		utils.Fail(c, 400, 40000, "参数解析失败")
		return
	}

	//新： 执行validator参数格式校验
	// 参数格式校验
	if err := utils.Validator.Struct(&req); err != nil {
		//旧：c.JSON(http.StatusBadRequest, gin.H{"msg": "参数格式不合法"})
		utils.FailWithValidator(c, err)
		return
	}

	err := service.UpdateUser(&req)
	if err != nil {
		//旧：c.JSON(http.StatusInternalServerError, gin.H{"msg": "更新失败:" + err.Error()})
		utils.Fail(c, 500, 50000, "更新失败:"+err.Error())
		return
	}
	//旧：c.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
	utils.Success(c, nil)
}

// DeleteUser 删除用户接口
func DeleteUser(c *gin.Context) {
	// 用DTO接收路径参数
	var param dto.UserIdUri
	if err := c.ShouldBindUri(&param); err != nil {
		//旧：c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		utils.Fail(c, 400, 40000, "参数解析失败")
		return
	}

	//新： 执行validator参数格式校验
	if err := utils.Validator.Struct(&param); err != nil {
		//旧：c.JSON(http.StatusBadRequest, gin.H{"msg": "参数格式不合法"})
		utils.FailWithValidator(c, err)
		return
	}

	err := service.DeleteUser(param.ID)
	if err != nil {
		//旧：c.JSON(http.StatusInternalServerError, gin.H{"msg": "删除失败:" + err.Error()})
		utils.Fail(c, 500, 50000, "删除失败:"+err.Error())
		return
	}
	//旧：c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
	utils.Success(c, nil)
}

// 创建时间：2026/8/17 下午7:01
package handler

import (
	"learn_layer_demo/model"
	"learn_layer_demo/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserById  查询用户接口
func GetUserById(c *gin.Context) { //在执行/:id查询指定用户时，自动触发该函数，c中携带了前端传来的数据
	type UriParam struct {
		ID uint `uri:"id"`
	}
	var param UriParam
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	id := param.ID
	user, err := service.GetUserInfo(id) //将id传给service层的查询函数
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "查询失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// CreateUser 新增用户接口
func CreateUser(c *gin.Context) {
	var user model.User
	// 绑定前端传来的json
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数解析失败"})
		return
	}
	err := service.AddUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "新增失败:" + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "新增成功", "data": user})
}

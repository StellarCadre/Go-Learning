// 创建时间：2026/8/18 下午8:19
// create:2026‑08‑18
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response 统一接口返回结构体
// 约定全部接口都返回这套json结构
type Response struct {
	Code int         `json:"code"` // 业务码：0代表成功；非0代表业务错误
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 返回数据，没有数据返回nil
}

/*
0 = 业务处理成功
40000 = 参数类错误
50000 = 业务 / 服务内部错误
*/

// Success 成功响应
// c:gin上下文；data:需要返回的数据
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "操作成功",
		Data: data,
	})
}

// Fail 失败响应
// httpCode http状态码；bizCode业务错误码；msg错误提示
/*
httpCode：给浏览器、网络层面看；例如参数错返回 400，资源找不到返回 404。
bizCode：给前端程序做逻辑判断，例如前端判断如果 code !=0 就弹出 msg 提示。
*/
func Fail(c *gin.Context, httpCode int, bizCode int, msg string) {
	c.JSON(httpCode, Response{
		Code: bizCode,
		Msg:  msg,
		Data: nil,
	})
}

// FailWithValidator 处理validator参数校验错误，解析出具体字段错误信息
func FailWithValidator(c *gin.Context, err error) {
	// 类型断言，判断是不是validator校验错误
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 如果不是validator错误，返回通用提示
		Fail(c, http.StatusBadRequest, 40000, "参数格式不合法")
		return
	}

	// 取第一条错误信息返回；也可以循环拼接全部错误
	/*
		fe.Field()：出问题的结构体字段名字，例如 Name、ID、Age
		fe.Tag()：触发失败的校验标签名字字符串，"required" / "gt" / "gte" 这些
		fe.Param()：标签后面写的参数，比如 gt=0，参数就是字符串 "0"
	*/
	/*
		tag标签	     含义
		required	 必填
		gt	         大于
		gte	        大于等于
		lte	        小于等于
		min	        最小长度
		max	        最大长度
	*/
	fe := errs[0] // 拿到第一条校验错误对象
	msg := ""     // 定义空字符串，用来存要返回给前端的中文提示
	// switch 判断 fe.Tag() 的值，也就是看是哪一条validate规则报错了
	switch fe.Tag() {
	case "required":
		// 如果Tag是required → 代表【必填】校验失败
		msg = fe.Field() + "为必填项"
	case "gt":
		// gt = greater than，必须大于；比如 validate:"gt=0"
		msg = fe.Field() + "必须大于" + fe.Param() //// msg = "ID必须大于0"
	case "gte":
		// gte = greater than or equal，大于等于；比如 validate:"gte=0"
		msg = fe.Field() + "必须大于等于" + fe.Param()
	case "lte":
		// lte = less than or equal，小于等于；比如 validate:"lte=150"
		msg = fe.Field() + "必须小于等于" + fe.Param()
	case "min":
		// min 最小长度，validate:"min=1"
		msg = fe.Field() + "长度不能小于" + fe.Param()
	case "max":
		// max 最大长度，validate:"max=20"
		msg = fe.Field() + "长度不能大于" + fe.Param()
	default:
		// 上面所有case都没有匹配上，走这里兜底
		msg = fe.Field() + "参数非法"
	}
	Fail(c, http.StatusBadRequest, 40000, msg)
}

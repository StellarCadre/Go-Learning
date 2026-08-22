// 创建时间：2026/8/18 下午7:49

package utils

import "github.com/go-playground/validator/v10"

// Validator 全局校验器单例，整个项目复用同一个实例，不要每次请求new，性能更好
var Validator *validator.Validate

func init() {
	// init包初始化函数，程序启动自动执行，初始化校验器
	Validator = validator.New()
}

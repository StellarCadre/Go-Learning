// 创建时间：2026/5/27 下午7:48
package main
/*
把函数当指令存起来（菜单 / 系统操作）
把函数当参数传递（通用遍历 / 处理逻辑）
包装函数，加统一功能（日志 / 耗时 / 校验）
返回一个新函数（动态生成工具）
 */
import "fmt"

////高阶函数：把函数作为参数传递给另一个函数或其他类型变量
func login(){
	println("登录成功")
}
func logout(){
	println("登出成功")
}
func register(){
	println("注册成功")
}

func main() {
	// 🔥 这里就是 【高阶函数】 的体现！
	// ==============================
	// map 的 key 是 int
	// map 的 value 是 【函数类型】 func()
	// 把函数 login / register / logout 当成值存进 map！
	var funMap=map[int]func(){
		1: login,
		2: register,
		3: logout,
	}
   fmt.Println("欢迎进入该系统，请选择操作：1.登录 2.注册 3.退出")
   var operation int
   fmt.Scan(&operation)
	// 从 map 中根据用户输入的数字取出对应的函数
	// fun 拿到的是：函数本身
	// ok 拿到的是：是否找到这个key
   fun, ok := funMap[operation]  //fun：匹配到的函数（若存在） ok：布尔值，true表示 key 存在，false表示不存在
   if ok {
	   fun()
   }else {
	   fmt.Println("输入错误")
   }
}

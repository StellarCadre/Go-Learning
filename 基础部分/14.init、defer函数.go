// 创建时间：2026/5/27 下午9:19
package main

import "fmt"

/*
init函数：在main函数执行之前执行（自动顺序进行），复杂变量的初始化、资源的申请、运行前检查、组件与驱动的“隐式注册”、运行环境的校验与准备。
         除了下面这种将init函数写在main函数外，还可以在自定义文件中定义init函数，然后使用import导入：import 包名。
         注意，一个包中，只有首字母大写的变量、函数、方法、属性才能被其他包访问。
         注意，导入的包，如果没有被使用，会报错（不过现在的编译器会自动帮忙删除无用包），或者下划线_来引用包，这样就不会报错。如：_ "fmt"。
         注意，导入的包也可以重命名，如：import fmt1 "fmt"
*/

/*
defer函数：这些调用直到函数结束前或return后才被执行，用于释放资源，如最后关闭文件，关闭连接。多个defer语句，按先进后出的顺序执行，谁离return近，谁先执行。
          defer 函数名()，这个函数名可以是普通函数，也可以是匿名函数。

*/

var user string = "Aurora"
var password string = "123456"

func init() {
	println("init1")
	user = "linhuahua"
}
func init() {
	println("init2")
	password = "00000000"
}

func Test() {
	fmt.Println("nihao")
}

func main() {
	println("main")
	fmt.Println(user, password)

	defer func() { //defer 匿名函数()
		println("defer1")
	}()

	defer Test() //defer 普通函数名()

	defer func() { //等价与defer fmt.Println("defer2")
		println("defer2")
	}()

	return
}

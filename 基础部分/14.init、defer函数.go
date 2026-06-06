// 创建时间：2026/5/27 下午9:19
package main

import "fmt"

//init函数：在main函数执行之前执行（顺序进行），用于初始化。
//defer函数：这些调用直到return前才被执行，用于释放资源。多个defer语句，按先进后出的顺序执行，谁离return近，谁先执行。

var user string= "Aurora"
var password string= "123456"
func init(){
	println("init1")
	user="linhuahua"
}
func init(){
	println("init2")
	password="00000000"
}


func main() {
    println("main")
	fmt.Println(user,password)



	defer func(){
		println("defer1")
	}()
	defer func(){  //等价与defer fmt.Println("defer2")
		println("defer2")
	}()
	return
}

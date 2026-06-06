// 创建时间：2026/5/27 下午8:57
package main

func copyName(name string){  //这里的name与来源处的name只有简单的赋值关系，地址并不一样
	name="linhuahua"
	println(name)

}

func changeName(name *string){  //这里将来源处的name的地址传递给name，两个变量的地址是一致的
	//下面尝试将name的值进行修改
	//name="gaoqiqiang"  //必定报错，因为这里name是用来存来源处地址的，不能直接赋值
	*name = "gaoqiqiang"  //需要使用*来解引用，这样才能修改name的值

}

func main() {
    var name="Aurora "
	println(name)
	copyName(name)  //值传递，传递的是name值的拷贝，所以这里name的值不会改变
	println(name)
	changeName(&name) //引用传递，传递的是name地址，所以这里name的值会改变
	println(name)
}

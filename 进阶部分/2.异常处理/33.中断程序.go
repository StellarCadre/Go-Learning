// 创建时间：2026/6/2 下午10:24
package main

import (
	"fmt"
	"os"
)

/*
对于中断，一般出现在初始化函数init中，如进行文件读取、等。
当出现错误时，不是进行向上抛，而是直接终止.
可使用log.Fatalln()
或者
panic

 */


func init(){
	_,err:=os.Open("./test.txt")
	if err!=nil{
		//log.Fatalln("打开文件失败：",err) //直接终止程序，并打印：2026/06/03 15:58:43 打开文件失败： open ./test.txt: The system cannot find the file specified.
		//或者使用panic(err)
		panic("打开文件失败") //直接终止程序，并打印：出错的堆栈信息
	}
}


func main() {
    fmt.Println("main")
}

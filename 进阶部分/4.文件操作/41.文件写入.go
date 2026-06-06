// 创建时间：2026/6/3 下午8:46
package main

import (
	"fmt"
	"os"
)

/*
文件操作中，常见的一些模式（可通过位或操作符|进行组合）：
O_RDONLY  : 以只读模式打开文件
O_WRONLY  : 以只写模式打开文件
O_RDWR    : 以读写模式打开文件
O_APPEND  : 以追加方式写入到文件末尾（不会覆盖原有内容）
O_CREATE  : 如果文件不存在则创建文件；若已存在则直接打开
O_TRUNC   : 打开文件时清空原有内容（仅对可写文件有效）
O_SYNC    : 同步写入，确保数据立即刷到磁盘（而非缓存）
O_EXCL    : 与O_CREATE配合使用，若文件已存在则报错（避免覆盖）

文件权限说明（第三个参数）只用于Linux中：
0666表示：所有者、同组用户、其他用户均拥有读和写权限（rw-rw-rw-）
- 第一位0表示八进制数标识
- 第二位6（4+2）：所有者权限（读+写）
- 第三位6（4+2）：同组用户权限（读+写）
- 第四位6（4+2）：其他用户权限（读+写）
注：Windows系统下文件权限会被映射，此参数仅作兼容标识
*/

func main() {
	//第一个参数是文件路径，第二个参数是文件打开模式，第三个参数是文件权限
    file,err:=os.OpenFile("C:\\Users\\Aurora\\Desktop\\Go Project\\进阶部分\\4.文件操作\\Test.txt",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)  //os.OpenFile可读、可写、可创建、可追加、可覆盖
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	n,err:=file.WriteString("马尾卡，我车坏了")
	if err!=nil{
		panic(err)
	}
	fmt.Println(n)
}

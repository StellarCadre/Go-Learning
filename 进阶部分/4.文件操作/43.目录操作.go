// 创建时间：2026/6/3 下午9:56
package main

import (
	"fmt"
	"os"
)

func main() {
    dir,err:=os.ReadDir("C:\\Users\\Aurora\\Desktop\\Go Project\\进阶部分\\4.文件操作") //dir中存放有很多内容
	if err!=nil{
		fmt.Println("读取目录失败")
	}

	for _,entry:=range dir{
		information,err:=entry.Info()
		if err!=nil{
			fmt.Println("获取文件信息失败")
		}
		// 打印当前条目的核心信息
		// 输出内容：名称 | 是否为目录 | 文件大小（字节） | 最后修改时间
		fmt.Println(entry.Name(),entry.IsDir(),information.Size(),information.ModTime())
	}
}

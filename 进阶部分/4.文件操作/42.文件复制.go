// 创建时间：2026/6/3 下午9:37
package main

import (
    "io"
	"log"
	"os"
)

func main() {
    file1,err:=os.Open("C:\\Users\\Aurora\\Desktop\\Go Project\\进阶部分\\4.文件操作\\成绩.jpg") //打开待复制文件
    if err!=nil{
		log.Fatalf("打开源文件失败：%v", err)
    }
    defer file1.Close()

    file2,err:=os.OpenFile("C:\\Users\\Aurora\\Desktop\\Go Project\\进阶部分\\4.文件操作\\成绩备份.jpg",os.O_CREATE|os.O_WRONLY,0666) //创建或打开目标文件
	if err!=nil{
		log.Fatalf("创建/打开目标文件失败：%v", err)
	}
	defer file2.Close()

	// 复制：源文件(file1) -> 目标文件(file2)
	written, err := io.Copy(file2, file1)
	if err != nil {
		log.Fatalf("文件复制失败：%v", err)
	}

	log.Printf("文件复制成功，共写入 %d 字节", written)


}

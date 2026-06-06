// 创建时间：2026/6/3 下午7:16
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	//一次性读取 os.ReadFile。若文件过大，会占用大量内存，不适合大文件读取。
    //byteData,err:=os.ReadFile("C:\\Users\\Aurora\\Desktop\\Go Project\\进阶部分\\4.文件操作\\斗破苍穹.txt")  //这样是绝对路径读取
	//byteData,err:=os.ReadFile("斗破苍穹.txt")  //这样是相对路径读取
	//byteData,err:=os.ReadFile("进阶部分/4.文件操作/斗破苍穹.txt")  //这样是相对于当前项目的路径读取
	//if err!=nil{
	//	panic(err)
	//}
	//fmt.Println(string(byteData))

	//分片读取，通常用于读取二进制文件
	//file,err:=os.Open("C:\\Users\\Aurora\\Desktop\\Go Project\\进阶部分\\4.文件操作\\斗破苍穹.txt")  // 1. 打开文件（只读模式）
	//if err!=nil{
	//	panic(err)   // 打开失败则终止程序（比如文件不存在）
	//}
    //defer file.Close()  // 2. 延迟关闭文件：不管后续代码如何执行，函数结束时一定会关闭文件
	//for{  // 3. 循环读取（直到文件末尾）
	//	var byteData =make([]byte,10)  // 3.1 定义缓冲区：每次最多读取10个字节（可根据需求调整，比如1024、4096）
	//	n, err := file.Read(byteData)  // 3.2 从刚才打开的file中读取数据到缓冲区byteData：若文件还有数据：n=实际读取的字节数。// 若文件已读完：n=0，err=io.EOF
	//	if err ==io.EOF {   // 3.3 终止循环：读取到文件末尾则退出
	//		break
	//	}
    //    fmt.Println(string(byteData),n) // 3.4 打印本次读取的内容和实际字节数
	//}

	//按行读取，需要用到缓存。
	//file,err:=os.Open("C:\\Users\\Aurora\\Desktop\\Go Project\\进阶部分\\4.文件操作\\斗破苍穹.txt")
	//if err!=nil{
	//	panic(err)
	//}
    //defer file.Close()
	//buffer:=bufio.NewReader(file) //这里的buffer是缓存，可以理解为一个容器，用于存放读取的数据。
	//for  {
	//	line,_,err:=buffer.ReadLine() //这里的line是读取的一行数据，省略的是这一行是否太长、没读完？err是错误信息，如果读取到文件末尾，err是io.EOF
	//	if err!=nil{
	//		break
	//	}
	//	fmt.Println(string(line),err)
	//}

	//指定分隔符读取（默认是按行读取）
	file,err:=os.Open("C:\\Users\\Aurora\\Desktop\\Go Project\\进阶部分\\4.文件操作\\斗破苍穹.txt")  //os.open 只能是打开文件
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	buffer:=bufio.NewScanner(file)  //创建一个 “扫描器”

	buffer.Split(bufio.ScanWords)  // 设置按【空格】分割，而不是默认按行
	//buffer.Split(bufio.ScanLines)  // 默认就是这个按行
	//buffer.Split(bufio.ScanRunes)  // 按字符
	//buffer.Split(bufio.ScanBytes)  //按字节
	for buffer.Scan(){  //有数据返回true，for成立
		fmt.Println(buffer.Text())  // 获取读到的内容
	}




}

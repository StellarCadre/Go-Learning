// 创建时间：2026/6/2 下午9:17
package main

import (
	"errors"
	"fmt"
)

//Go 的错误向上抛 = 一层接一层手动传！
//这里用到的error是来处理普通错误的。

//函数一开头声明了什么返回值，空 return 就自动返回什么。

func divide(a,b int)(result int,err error){
	 if b==0{
		 err=errors.New("除数不能为0")
		 return  // 直接返回result, err：result默认是0，err是上面创建的错误
	 }
	 result=a/b
	 return   // 直接返回result, err。result是计算结果，err默认是nil（空，代表没错误）
}

func Test()(result int,err error){
    res,err:=divide(1,0)
	if err!=nil{
		return  //直接返回result, err，并把错误“抛”给调用Test的地方（也就是main函数）
	}

	res=res+1  //无错误时，可以选择性继续处理
	return
}


func main() {

     result,err:=Test()
	 if err!=nil{  // 如果Test返回了错误
		 fmt.Println("错误：",err)
		 return  //这里的return什么都不返回，而是用于main终止。
	 }
	 fmt.Println(result)


}

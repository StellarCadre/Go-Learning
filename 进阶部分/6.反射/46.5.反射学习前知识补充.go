// 创建时间：2026/6/20 下午7:31
/*
变量：type+value，总称pair
其中，type：
    static type静态类型，如int、string
    concrete具体类型，interface所指向的类型
一个变量的pair，其被传递到其他变量（赋值）时，其他变量也拿到了该pair。

*/

package main

func main() {
	//例子1
	//var a string
	////这里a的本质：pair<static type:string,value:"abc">
	//a="abc"
	//// 空接口 interface{} 底层结构：(动态类型 _type, 动态值 data)
	//var allType interface{}  //allType的动态类型是interface，内部存储的type是动态类型
	//allType = a  //这里allType的本质：pair<type:string,value:"abc">
	//str,_ :=allType.(string)
	//fmt.Println(str)

	//例子2
	////tty的本质：pair<type:*os.File,value:"dev/tty"文件描述符>
	//tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	////r的本质：pair<type:  ,value:  >
	//var r io.Reader
	////赋值之后，r的本质：pair<type:*os.File  ,value:"dev/tty"文件描述符  >
	//r = tty
	//
	////w的本质：pair<type:  ,value:  >
	//var w io.Writer
	////赋值之后，w的本质：pair<type:*os.File  ,value:"dev/tty"文件描述符  >
	//w = r.(io.Writer)
}

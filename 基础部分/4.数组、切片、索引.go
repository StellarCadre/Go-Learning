package main

import (
	"fmt"
	"sort"
)

/*
数组和切片在定义时的区别：
数组：定义时必须指定长度，且长度不能改变。var idList2 [5]int
切片：定义时不需要指定长度，但长度可以改变。var idList2 []int

*/

func test1(nameList []string) {}

func main() {
	//数组的定义和索引
	//数组的标准定义形式为： var 数组名 [长度]类型 = [长度]类型{植}，其中前面的[长度]类型可以省略
	var idList1 [5]int = [5]int{1001, 1002, 1003, 1004, 1005}
	fmt.Println(idList1)
	fmt.Println(idList1[4]) //数组下标从0开始，值分别为[0] [1] [2] [3] [4]共5个。
	fmt.Println(idList1[len(idList1)-2])
	//数组的遍历方法：使用range关键字,能够自动返回内容。
	for index, value := range idList1 { //index是索引，value是值
		fmt.Println(index, value)
	}
	//打印数组的数据类型：
	fmt.Printf("idList1的数据类型是：%T\n", idList1) //类型是[5]int,那么其在任何处理过程中都不会改变，长度也是无法改变
	//数组另一种定义形式：
	var idList2 [5]int //只定义但不初始化，值全为0
	for i := 0; i < len(idList2); i++ {
		fmt.Println(idList2[i])
	}

	fmt.Println("-----------------------下面来讲数组的升级版：切片，又名动态数组----------------------------------")
	//切片slice：数组的升级版，来解决定义好的数组长度被定死的问题，后期无法修改的问题(在python中叫列表list)。又称动态数组
	//切片的定义和索引，定义时只需要写[]类型即可，不需要指定长度。
	var nameList []string
	nameList = append(nameList, "空山新雨后，天气晚来秋。") //格式为： 切片名 = append(切片名, 值)
	nameList = append(nameList, "明月松间照，清泉石上流。")
	fmt.Println(nameList)
	fmt.Println(nameList[1]) //切片的索引和数组一样，从0开始
	test1(nameList)          //切片传给其他函数时，为引用传递（传递切片的地址），所以在函数中对切片的修改会影响到原切片的值。但数组是值传递，在函数中对数组的修改不会影响到原数组的值。

	//切片只定义不初始化的情况（只声明，没赋值，没分配内存）：
	var nameList1 []string
	//fmt.Println(nameList1[0])     //切片未初始化，无法使用索引，会报错
	fmt.Println(nameList1 == nil) //切片未初始化，值为nil,即空nil

	//切片定义并初始化
	var nameList11 []int = []int{1, 2, 3, 4, 5}
	fmt.Println(nameList11)

	//切片初始化为空切片的方法（底层分配了空数组内存）：
	var nameList2 []string = []string{} // 方法一：显式声明类型并初始化为空切片
	var nameList3 = []string{}          // 方法二：类型推断并初始化为空切片
	nameList4 := []string{}             // 方法三：短变量声明（推荐）
	nameList5 := make([]string, 0)      // 方法四：用 make 创建（推荐，可同时指定容量0）
	fmt.Println(nameList2 == nil)       // false   // 验证：初始化后的切片都不是 nil
	fmt.Println(nameList3 == nil)       // false
	fmt.Println(nameList4 == nil)       // false
	fmt.Println(nameList5 == nil)       // false

	//使用make函数创建指定长度、容量和类型的切片。长度是切片的长度。容量一般等于长度（用的少）
	IDList1 := make([]int, 5) //长度为5的切片。全是0。不写容量，默认和长度一样。
	//IDList1 := make([]int, 5,7) //长度为5，容量为7的切片。5长度的意思是当前有5个有意义的值，7容量的意思是最多存放7个值。容量不小于长度。5个0，剩余2个无法访问，除非使用append()方法在切片末尾添加元素。
	fmt.Println(IDList1)          //int默认值是0
	IDList1[0] = 666              //在索引0的位置上赋值为666
	IDList1 = append(IDList1, 21) //append 是在切片末尾添加元素，不会覆盖原来的内容，会使长度+1。 不用管容量，容量不够会自动扩容。
	IDList1 = append(IDList1, 55)
	fmt.Println(IDList1)
	fmt.Printf("IDList1的长度%v,容量%v\n", len(IDList1), cap(IDList1))

	//切片的截取,左闭右开
	IntList := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(IntList[2:5]) //从索引2开始，到索引5结束，不包含索引5。结果：[3 4 5]
	fmt.Println(IntList[2:])  //从索引2开始，到结束。结果：[3 4 5 6 7 8 9 10]
	fmt.Println(IntList[:5])  //从开始，到索引5结束，不包含索引5。结果：[1 2 3 4 5]
	fmt.Println(IntList[:])   //从开始，到结束。结果：[1 2 3 4 5 6 7 8 9 10]

	//切片的浅复制和深复制copy
	IntList1 := IntList[:] //切片的复制，复制的是地址，不是值。
	IntList1[0] = 1000     //故对切出来的部分赋值，原切片也会被修改。这叫浅复制，深复制需要使用copy()函数。
	fmt.Println(IntList)
	fmt.Println(IntList1)
	IntList2 := make([]int, 5)
	copy(IntList2, IntList) //这样得到的是深复制，不会影响原切片。

	//切片的排序方法
	IDList2 := []int{4, 3, 5, 1, 2}
	fmt.Println("默认输出：", IDList2)
	sort.Ints(IDList2) //sort.Ints是int型元素升序排序方法，等同：sort.Sort(sort.IntSlice(IDList2))
	fmt.Println("升序输出：", IDList2)
	sort.Sort(sort.Reverse(sort.IntSlice(IDList2))) //Go中没有降序的方法，只能使用sort.Sort()万能排序方法。即：sort.Sort(sort.Reverse(sort.IntSlice(IDList2)))
	fmt.Println("降序输出：", IDList2)
}

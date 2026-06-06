package main

import (
	"fmt"
	"sort"
)

func main() {
	//数组的定义和索引
	var idList [5]int = [5]int{1001, 1002, 1003, 1004, 1005} //数组的标准定义形式为： var 数组名 [长度]类型 = [长度]类型{植}，其中前面的[长度]类型可以省略
	fmt.Println(idList)
	fmt.Println(idList[4]) //数组下标从0开始，值分别为[0] [1] [2] [3] [4]共5个。
	fmt.Println(idList[len(idList)-2])

	//数组的升级版：切片slice，来解决定义好的数组长度被定死的问题，后期无法修改的问题(在python中叫列表list)
	var nameList []string
	nameList = append(nameList, "空山新雨后，天气晚来秋。") //格式为： 切片名 = append(切片名, 值)
	nameList = append(nameList, "明月松间照，清泉石上流。")
	fmt.Println(nameList)
	fmt.Println(nameList[1]) //切片的索引和数组一样，从0开始
	fmt.Println(nameList[0])

	//切片只定义不初始化的情况（只声明，没分配内存）：
	var nameList1 []string
	//fmt.Println(nameList1[0])     //切片未初始化，无法使用索引，会报错
	fmt.Println(nameList1 == nil) //切片未初始化，值为nil,即空nil
	// 切片的初始化为空切片的方法（底层分配了空数组内存）：
	var nameList2 []string = []string{} // 方法一：显式声明类型并初始化
	var nameList3 = []string{}          // 方法二：类型推断
	nameList4 := []string{}             // 方法三：短变量声明（推荐）
	nameList5 := make([]string, 0)      // 方法四：用 make 创建（推荐，可同时指定容量）
	fmt.Println(nameList2 == nil)       // false   // 验证：初始化后的切片都不是 nil
	fmt.Println(nameList3 == nil)       // false
	fmt.Println(nameList4 == nil)       // false
	fmt.Println(nameList5 == nil)       // false

	//使用make函数创建指定长度、容量和类型的切片。长度是切片的长度。容量一般等于长度（用的少）
	IDList1 := make([]int, 5)     //5个元素的切片，长度和容量都是5。
	fmt.Println(IDList1)          //int默认值是0
	IDList1[0] = 666              //在索引0的位置上添加元素666
	IDList1 = append(IDList1, 21) //append 是在切片末尾添加元素，不会覆盖原来的内容，会使长度+1
	IDList1 = append(IDList1, 55)
	fmt.Println(IDList1)

	//切片的排序方法
	IDList2 := []int{4, 3, 5, 1, 2}
	fmt.Println("默认输出：", IDList2)
	sort.Ints(IDList2) //sort.Ints是int型元素升序排序方法，等同：sort.Sort(sort.IntSlice(IDList2))
	fmt.Println("升序输出：", IDList2)
	sort.Sort(sort.Reverse(sort.IntSlice(IDList2))) //Go中没有降序的方法，只能使用sort.Sort()万能排序方法。即：sort.Sort(sort.Reverse(sort.IntSlice(IDList2)))
	fmt.Println("降序输出：", IDList2)
}

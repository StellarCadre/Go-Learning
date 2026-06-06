// 创建时间：2026/6/3 下午7:08
package main

func main() {
    type Myslice1 []int  //不使用泛型时定义切片
	var myslice1 = Myslice1{1,2,3}

	type Myslice2[T int|string] []T //使用泛型时定义切片
	var myslice1 = Myslice2[int]{1,2,3}
	var myslice2 = Myslice2[string]{"a","b","c"}
}

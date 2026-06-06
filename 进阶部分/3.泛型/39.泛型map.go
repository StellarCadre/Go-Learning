// 创建时间：2026/6/3 下午7:11
package main

func main() {
	type Mymap1 map[string]int  ////不使用泛型时定义map
	mymap1:=Mymap1{
		"a":1,
		"b":2,
		"c":3,
	}

    type Mymap2[T int|string,K int|string] map[T]K
	Mymap2:=Mymap2[string,int]{
		"a":1,
		"b":2,
		"c":3,
	}
}


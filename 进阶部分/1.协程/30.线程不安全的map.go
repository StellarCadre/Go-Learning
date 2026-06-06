// 创建时间：2026/6/2 下午8:49
package main

import (
	"fmt"
)

//下面这样的代码会报错：fatal error: concurrent map read and map write，因为是并发读写map

func main() {
	var maps = map[int]string{}
    go func() {  //对map进行写操作
		for{
			maps[1] = "a"
		}
	}()

	go func() {  //对map进行读操作
		for{
			fmt.Println(maps[1])
		}
	}()

	select {}

}

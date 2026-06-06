// 创建时间：2026/6/2 下午8:59

package main

import (
	"fmt"
	"sync"
)

// 线程安全的map要将map声明为sync.Map类型，而不是普通的map[int]string类型。sync.Map源码中默认使用了读写锁，因此可以支持并发读写。
// 同时读写map的操作也需要使用sync.Map提供的Load和Store方法，而不是直接使用map[key]=value的方式。

func main() {
	var maps = sync.Map{}
	go func() {  //对map进行写操作
		for{
			maps.Store(1, "a")
		}
	}()

	go func() {  //对map进行读操作
		for{
			val,ok:=maps.Load(1)
			fmt.Println(val,ok)
		}
	}()

	select {}

}


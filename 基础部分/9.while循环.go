// 创建时间：2026/5/25 下午8:17
package main

func main() {
	//go 没有while 循环,可以使用for改造得到
	var sum int
	var i int = 1
	for i <= 100 { //while先判断，再执行。for是先执行，再判断
		sum += i
		i++
	}

	var add int
	var n int = 1
	for { //do while先执行，再判断
		add += n
		n++
		if n > 100 {
			break
		}
	}

	//切片的遍历
	var list = []string{"你", "我", "他"}
	for i := 0; i < len(list); i++ { //用i来记录下标。切片的长度是len(list)，下标范围是0-len(list)-1。这是遍历的方法一。
		println(i, list[i])
	}
	for index, value := range list {
		println(index, value)
	}

}

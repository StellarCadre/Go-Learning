package main

import "fmt"

func main() {

	//求和
	var number int
	var sum int
	for number = 1; number <= 10; number++ { //第一个是初始化语句，第二个是条件语句，第三个是递增语句
		fmt.Print(number, " ")
		sum += number
	}
	fmt.Println("\n1到10的和：", sum)

	//死循环
	//for {
	//	fmt.Println(time.Now())
	//	time.Sleep(1 * time.Second)
	//}

	//break和continue的使用,break用于跳出循环，continue用于跳过本次循环
    //99乘法表
	fmt.Println("\n99乘法表：")
    for i:=1;i<=9;i++{
		for j:=1;j<=i;j++{
			fmt.Printf("%d*%d=%d\t", j, i, j*i)
		}
		fmt.Println()
	}

	fmt.Println("\n100以内的奇数求和：")
	var sum1 int
	for i:=1;i<=100;i++{
		if i%2==0 {
			continue
		}
		sum1 += i
	}
	fmt.Println(sum1)


	fmt.Println("\n当i==5时，结束循环")
	for i:=0;i<=10;i++{
		if i==5{
			break
		}
		fmt.Println(i)

	}


}

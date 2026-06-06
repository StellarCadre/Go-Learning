package main

import "fmt"

func main() {

	//整数类型
	var a1 = 1            //默认是int类型
	var a2 int = 2        //显式声明为int类型
	var a3 uint8 = 255    //前面的u表示无符号的int型（只能保存正值），后面的8表示占8位，范围是0-255，即00000000-11111111
	var a4 int8 = -128    //8表示占8位，范围是-128-127
	var a5 uint16 = 32767 //16表示占16位，范围是0-65535,即从0到2的16次方减1
	var a6 int16 = -32768 //16表示占16位，范围是-2的15的次方到2的15次方减1

	//浮点型：float32和float64
	var b1 float32 = 3.1415926 //范围是3.4e-38到3.4e+38
	var b2 float64 = 3.1415926 //范围是1.8e-308到1.8e+308，未显式声明时默认是float64

	fmt.Println(a1, a2, a3, a4, a5, a6, b1, b2)

	//字符型：单字节字符byte
	var c1 byte = 'a'
	fmt.Println(c1)        //在Go中，字符本质是整数，所以直接输出时只会是该字符对应的ASCII码值97
	fmt.Printf("%c\n", c1) //使用%c才能表示原本的字符
	var c2 uint = 97
	fmt.Println(c2)        //在Go中，数字就是数字，所以直接输出时只会是该数字的值
	fmt.Printf("%c\n", c2) //使用%c时，会把该数字转换为对应的ASCII码值，然后输出。综上所述，在int8范围内，对整数使用%c和对字符使用%d是高度相关的

	//字符型：多字节字符rune，如中文、日文、韩文等
	var d1 rune = '中'
	fmt.Println(d1)        //在Go中，字符本质是整数，所以直接输出时只会是该字符对应的Unicode码值20013
	fmt.Printf("%c\n", d1) //使用%c才能表示原本的字符

	//字符串型：string
	var e1 string = "hello Go"
	fmt.Println(e1)
	fmt.Println("转义字符的使用：------------------------------------")
	fmt.Println("1:北国\t风光") //\t表示制表符，即一个制表符的宽度是8个空格
	fmt.Println("2:北国\n风光") //\n表示换行，即换行符
	fmt.Println("3:北国\r风光") //\r表示回车，即回到当前行的开头，相当于删除当前行，然后从头开始写
	fmt.Println("4:北国\b风光") //\b表示退格，即删除一个字符，相当于把光标往左移动一个字符的位置
	fmt.Println("北国'的'风光")  //在双引号中直接加单引号可行
	//fmt.Println("北国"的"风光") //在双引号中直接加双引号不可行，需要用到转义字符\
	fmt.Println("北国\"的\"风光") //在双引号中直接加双引号不可行，需要用到转义字符\
	fmt.Println("多行字符串的使用：------------------------------------")
	fmt.Println(`
		北国有佳人，
		绝世而独立。
		一顾倾城色，
		再顾倾国色。
		人间车马多如簇星，
		但见京华如梦里。
		不羡清华北大，
		只羡美人如月。
	`) //多行字符串的使用，使用反引号``，在其中可以直接换行，不受制表符、换行符、回车符的影响。反引号在esc的正下方

	//布尔类型：bool:true和false,无法参与数值运算，也无法进行类型转换
	var f1 bool = true
	fmt.Println(f1)
	var f2 bool = false
	fmt.Println(f2)
	fmt.Println(!f2)

	//零值问题：只对基本数据类型进行定义但不赋值时，会自动赋予默认值。
	var g1 int
	var g2 float32
	var g3 byte
	var g4 rune
	var g5 string
	var g6 bool
	fmt.Printf("%#v\n", g1) //默认值是0
	fmt.Printf("%#v\n", g2) //默认值是0
	fmt.Printf("%#v\n", g3) //默认值是0
	fmt.Printf("%#v\n", g4) //默认值是0
	fmt.Printf("%#v\n", g5) //默认值是""
	fmt.Printf("%#v\n", g6) //默认值是false
}

// 创建时间：2026/5/31 下午10:18
package main

import "fmt"

/*
接口是一组仅包含方法名、参数、返回值的方法的集合，未实现具体方法
可以解决相似功能（仅传入对象不同 ）的方法冗余的问题，如下：
*/

type AnimalInterface interface {  //// 定义一个接口：约定只要有Sing()方法的类型，都属于这个接口，正好Chicken和Tiger都符合这个接口
	Sing()
	GetName() string
}

type Chicken struct {
	name string
}
type Tiger struct {
	name string
}

func (c Chicken)Sing(){  //Chicken实现了Sing() → 自动属于AnimalInterface类型
	fmt.Println(c.name,"在唱歌")
}
func (t Tiger)Sing(){  // Tiger实现了Sing() → 也自动属于AnimalInterface类型
	fmt.Println(t.name,"在唱歌")
}
func (c Chicken)GetName() string {
	return c.name
}
func (t Tiger)GetName() string {
	return t.name
}

func sing1(a AnimalInterface){  // 任何实现了该接口的类型，都能传入这个函数
	a.Sing()
	fmt.Println(a.GetName())

	switch a.(type) {  //这叫断言，判断属于哪个类型
	case Chicken:
		fmt.Println("我是小鸡")
	case Tiger:
		fmt.Println("我是老虎")
	}
}

func main() {
	c1:=Chicken{name:"小鸡"}
	sing1(c1)
	d1:=Tiger{name:"小虎"}
	sing1(d1)

}

/*
对接口的理解：
接口里的方法 = 多个类型的「公共方法清单」，猫、狗、鸡、老虎 都会叫，接口就定义：Sing()，大家都实现这个方法就行。
像 sing1 (a 接口) 这样的函数，它不关心你传进来是什么具体类型，它只关心：你有没有实现接口里的方法（只要有，就能传进来），函数内部可以随意使用接口里定义的公共方法。所有实现了接口的类型，都能传入这个函数
空接口：interface{}  任何类型都实现了这个接口

 */

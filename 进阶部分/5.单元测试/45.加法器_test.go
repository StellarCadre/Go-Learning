// 创建时间：2026/6/4 下午4:46
package main

import "testing"

/*准则
测试文件必须以 _test.go 结尾：xxx_test.go
测试函数名必须：Test大写开头，TestAdd / TestSum，不能Testadd、testAdd
入参固定：(t *testing.T)，不能改参数
 */

//当前这样，一个测试函数只能测试一个案例
//func TestAdd(t *testing.T) {
//	result:=add(1,2)
//	if result!=3 {
//		t.Errorf("测试失败")
//		return
//	}
//	t.Logf("测试成功")
//}

//子测试:这样写的话，我根据名字想测哪个用例就测那个
//func TestAdd(t *testing.T) {
//	t.Run("gagf", func(t *testing.T) {
//		if  add(1, 2) !=3{
//			t.Errorf("测试失败")
//			return
//		}
//	})
//	t.Run("afsfsg", func(t *testing.T) {
//		if  add(1, 2) !=0{
//			t.Errorf("测试失败")
//			return
//		}
//	})
//	t.Logf("测试成功")
//}

//若测试用例很多，可以将这些统一放到一个cases中：
func TestAdd(t *testing.T) {
	cases:=[]struct{
		Name string
		A int
		B int
		Expected int
	}{
		{"a1",2,3,5},
		{"a2",3,5,8},
		{"a3",4,7,11},
	}
	for _, s := range cases {
		t.Run(s.Name, func(t *testing.T) {
			if add(s.A,s.B)!=s.Expected{
				t.Errorf("测试失败")
				return
			}
			t.Logf("测试成功")
		})
	}

}

/*
单元测试提供的日志方法：
Log ：打印日志，同时结束测试
Logf ：格式化打印日志，同时结束测试
Error ：打印错误日志，同时结束测试
Errorf ：格式化打印错误日志，同时结束测试
Fatal : 打印致命日志，同时结束测试
Fatalf : 格式化打印致命日志，同时结束测试
 */
// 创建时间：2026/6/4 下午10:49
package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Class struct {
	Name string `orm:"name"`
	Age  int    `orm:"age"`
}

func Find(obj any, query ...any) (sql string, err error) { //sql语句拼接函数.obj表示查询的结构体,query表示查询条件
	t := reflect.TypeOf(obj) //获取类型,在这里是Class结构体类型
	if t.Kind() != reflect.Struct {
		err = errors.New("obj不是结构体")
		return
	}
	var where string
	//接下来的验证条件的环节
	if len(query) > 0 {
		q := query[0]
		qs, ok := q.(string)
		if !ok {
			err = errors.New("query的第一个参数必须是字符串")
			return
		}
		if strings.Count(qs, "?")+1 != len(query) {
			err = errors.New("查询参数个数不匹配")
			return
		}

		//替换操作
		for _, a := range query[1:] {
			switch s := a.(type) {
			case string:
				qs = strings.Replace(qs, "?", fmt.Sprintf("'%s'", s), 1) //将实际的内容替换掉？占位符
			case int:
				qs = strings.Replace(qs, "?", fmt.Sprintf("'%d'", s), 1)
			}
		}
		where = "where" + qs
	}
	//拼接所有有orm的字段
	var columns []string
	v := reflect.ValueOf(obj)           //获取值
	for i := 0; i < v.NumField(); i++ { //使用一个for循环，来多次拿字段
		typename := t.Field(i)             //获取索引为i时的字段的名称、标签等。
		JsonTag := typename.Tag.Get("orm") //获取字段的标签中的json字符串，结果是name、age
		if JsonTag == "" {
			continue
		}
		columns = append(columns, JsonTag)
	}
	//获取结构体名字并小写
	name := strings.ToLower(t.Name()) + "s"

	//拼接sql语句
	sql = fmt.Sprintf("select %s from %s %s", strings.Join(columns, ","), name, where)
	return
}

func main() {
	//select name from class where name="aurora"
	sql1, err1 := Find(Class{}, "name=?", "aurora")
	fmt.Println(sql1, err1)

	//select name,age from class where name="aurora" and age=18
	sql2, err2 := Find(Class{}, "name=? and age=?", "aurora", 18)
	fmt.Println(sql2, err2)
	//select name,age from class
	sql3, err3 := Find(Class{})
	fmt.Println(sql3, err3)

}

// 创建时间：2026/6/23 下午8:37
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// 1. 发起HTTP GET请求，获取网页源码
	resp, err := http.Get("https://www.baidu.com")
	if err != nil {
		log.Fatal("请求网页失败：", err)
	}
	defer resp.Body.Close() // 必须关闭响应体，防止资源泄漏

	// 状态码非200直接报错
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("网页访问异常，状态码：%d", resp.StatusCode)
	}

	// 2. 加载HTML文档，生成goquery文档对象（对应jQuery的$()）
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("解析HTML失败：", err)
	}

	// 3. 基础选择器用法（完全对标jQuery选择器）
	fmt.Println("===== 页面标题 =====")
	title := doc.Find("title").Text()
	fmt.Println(title)

	fmt.Println("\n===== 页面所有超链接文本+地址 =====")
	// 遍历所有a标签
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		linkText := item.Text()          // 获取标签内文本
		href, exist := item.Attr("href") // 获取href属性
		if exist {
			fmt.Printf("第%d条：文本=%s | 链接=%s\n", index+1, linkText, href)
		}
	})

	// 进阶：class选择器、id选择器
	// doc.Find("#id名")  根据id筛选
	// doc.Find(".class名") 根据class筛选
	// doc.Find("div.container > p") 层级选择
}

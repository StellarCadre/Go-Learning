// main.go Gin 启动自动打开浏览器，无404
package main

//导入gin并使用
import (
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// openBrowser 跨平台打开浏览器
func openBrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	return exec.Command(cmd, args...).Start()
}

func main() {
	// 新建空白gin引擎，消除警告
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies([]string{})

	// 注册 /hello 接口
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello Go Modules 自动打开浏览器成功！")
	})

	addr := ":8080"
	// 关键：直接跳转到存在的 /hello 路由，不会404
	visitUrl := "http://localhost" + addr + "/hello"

	// 协程启动服务
	go func() {
		err := r.Run(addr)
		if err != nil {
			panic("端口占用：" + err.Error())
		}
	}()

	time.Sleep(500 * time.Millisecond)
	err := openBrowser(visitUrl)
	if err != nil {
		println("自动打开失败，请手动访问：", visitUrl)
	} else {
		println("✅ 已自动跳转浏览器：", visitUrl)
	}

	// 阻塞程序不退出
	select {}
}

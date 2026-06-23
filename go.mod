//这里放的是项目的依赖清单，等价于java的pom.xml，放在项目根目录，是整个 Go 模块的核心配置文件。

// 1. 模块路径：整个项目的唯一标识，导入包时必须用这个路径
//比如模块名是 Go_Project，项目内 utils 文件夹包，导入写法：import "Go_Project/utils"
module Go_Project

// 2. 约束本项目最低要求的Go语言版本
go 1.26

//3.require 包路径 版本声明依赖哪个第三方库、锁定精确版本，不会出现环境差异。
// (1). 直接依赖：你代码里手动import的第三方包，锁定固定版本
//require github.com/gin-gonic/gin v1.9.1

// (2). indirect 间接依赖：第三方包内部自己依赖的包，你代码没直接导入
//require github.com/go-playground/validator/v10 v10.15.0 // indirect

// (3). replace 依赖替换：把远程包换成本地修改后的源码，本地调试第三方库专用
//replace github.com/gin-gonic/gin => ./local-gin

// (4). exclude 版本排除：强制项目永远不用某个有bug的依赖版本
//exclude github.com/gin-gonic/gin v1.9.0

//4.日常操作自动维护
// 完全不用手动修改编辑，命令自动更新：
//新增 import 第三方包 → 执行go mod tidy，自动添加 require；
//升级 / 降级依赖 → go get 包名@版本，自动修改版本号；
//删除无用导入 → go mod tidy，自动清理废弃 require 行。
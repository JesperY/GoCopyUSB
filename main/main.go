package main

import (
	"GoCopyUSB/backend"
)

func main() {

	// 进行初始化操作，例如读取配置文件、检查权限，单进程运行等

	// 调用 listener 开始监听
	backend.Listener()

	// 错误处理？

}

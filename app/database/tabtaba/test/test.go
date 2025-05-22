package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
)

// 安装依赖后测试启动
func main() {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(800, 600, "Test Window", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	for !window.ShouldClose() {
		glfw.PollEvents()
	}
}

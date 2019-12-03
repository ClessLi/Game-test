package main

import (
	"github.com/ClessLi/Game-test/game"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"runtime"
)

const (
	width              = 800
	height             = 600
	WordWidth  float32 = 1500
	WordHeight float32 = 1000
)

var (
	windowName = "Test Game"
	game2D     = game.NewGame(width, height, WordWidth, WordHeight)
	deltaTime  = 0.0
	lastFrame  = 0.0
)

func main() {
	runtime.LockOSThread()
	window := initGlfw()
	defer glfw.Terminate()
	initOpenGL()
	game2D.Init()
	for !window.ShouldClose() {
		currentFrame := glfw.GetTime()
		deltaTime = currentFrame - lastFrame
		lastFrame = currentFrame

		glfw.PollEvents()
		game2D.ProcessInput(float32(deltaTime))
		game2D.Update(deltaTime)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		game2D.Render(deltaTime)
		window.SwapBuffers()
	}
}
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(width, height, windowName, nil, nil)
	if err != nil {
		panic(err)
	}
	window.SetKeyCallback(KeyCallback)

	window.MakeContextCurrent()
	return window
}
func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Viewport(0, 0, width, height)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.BLEND)
}
func KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Press:
		game2D.Keys[key] = true
	case glfw.Release:
		game2D.Keys[key] = false
	}
}

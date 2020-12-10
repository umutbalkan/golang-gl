package main

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	width  = 800
	height = 600

	vertexShaderSource = `
		#version 330
		layout (location = 0) in vec3 aPos;

		void main() {
			gl_Position = vec4(aPos, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 330
		out vec4 FragColor;

		uniform vec4 ourColor;

		void main() {
			FragColor = ourColor;
		}
	` + "\x00"
)

var (
	triangle = []float32{
		0, 0.3, 0,
		-0.3, -0.3, 0,
		0.3, -0.3, 0,
	}

	triangleBg = []float32{
		0, 0.32, 0,
		-0.32, -0.315, 0,
		0.32, -0.315, 0,
	}
)

// initGlfw initializes glfw and returns a Window to use.
func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "test-crossplatform", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func main() {
	runtime.LockOSThread()
	fmt.Printf("OS: %s, Architecture: %s\n", runtime.GOOS, runtime.GOARCH)
	window := initGlfw()
	program := initOpenGL()
	vao := makeVao(triangle)
	vao2 := makeVao(triangleBg)
	gl.ClearColor(0.3, 0.3, 0.3, 1.0)
	// render-loop
	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyEscape) == 1 {
			window.SetShouldClose(true)
		}
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// timer and uniform  stuff
		timeValue := glfw.GetTime()
		redValue := math.Sin(timeValue)/2 + 0.5
		gl.UseProgram(program)

		gl.BindVertexArray(vao2)
		vertexColorLocation2 := gl.GetUniformLocation(program, gl.Str("ourColor\x00"))
		gl.Uniform4f(vertexColorLocation2, float32(redValue/1.5), float32(redValue/2), float32(0.0), float32(1.0))
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangleBg)/3))

		vertexColorLocation := gl.GetUniformLocation(program, gl.Str("ourColor\x00"))
		gl.Uniform4f(vertexColorLocation, float32(redValue), float32(redValue), float32(0.0), float32(1.0))

		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

		glfw.PollEvents()
		window.SwapBuffers()
	}

	defer glfw.Terminate()
}

// makeVao initializes and returns a vertex array from the points provided.
func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// 4* because float32 is 4 bytes (use unsafe.sizeof maybe? seems unsafe)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

/*

 */

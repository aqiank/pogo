package pogo

import (
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/jackyb/go-sdl2/sdl_image"
)

const (
	FULLSCREEN = sdl.WINDOW_FULLSCREEN
	OPENGL = sdl.WINDOW_OPENGL
	HIDDEN = sdl.WINDOW_HIDDEN
	BORDERLESS = sdl.WINDOW_BORDERLESS
	RESIZABLE = sdl.WINDOW_RESIZABLE
	MINIMIZED = sdl.WINDOW_MINIMIZED
	MAXIMIZED = sdl.WINDOW_MAXIMIZED
	INPUT_GRABBED = sdl.WINDOW_INPUT_GRABBED
)

type PImage struct {
	texture *sdl.Texture
	Width, Height int
}

// internal variables
var window *sdl.Window
var renderer *sdl.Renderer
var strokeColor, fillColor uint32 = 0x000000, 0xffffffff
var strokeEnabled, fillEnabled = true, true
var frameRate = 60
var period = 1000 / frameRate

// public variables that users can use
var Width, Height int
var MouseX, MouseY int = 0, 0
var Key byte

// public assignable handler functions that users can use
var QuitHandler func() = nil
var KeyDownHandler, KeyUpHandler func() = nil, nil

// setup the window and renderer
func Setup(w, h int, flags uint32) {
	Width, Height = w, h
	window, renderer = sdl.CreateWindowAndRenderer(Width, Height, flags)
	Background(0xf0f0f0)
}

// assign a function that will be executed indefinitely when the program is running
func Loop(loopFunc func()) {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				goto out
			case *sdl.MouseMotionEvent:
				MouseX = int(t.X)
				MouseY = int(t.Y)
			case *sdl.KeyDownEvent:
				if KeyDownHandler != nil {
					KeyDownHandler()
				}
				Key = byte(t.Keysym.Sym)
			case *sdl.KeyUpEvent:
				if KeyUpHandler != nil {
					KeyUpHandler()
				}
			}
		}
		loopFunc()
		renderer.Present()
		sdl.Delay(uint32(period))
	}

out:
	if QuitHandler != nil {
		QuitHandler()
	}
}

// delay execution by specified amount of time (in milliseconds)
func Delay(millis int) {
	sdl.Delay(uint32(millis))
}

// get how many frames displayed per second
func Framerate() int {
	return frameRate
}

// set how many frames displayed per second
func SetFramerate(fps int) {
	frameRate = fps
	period = 1000 / frameRate
}

// fill solid color onto the whole screen
func Background(color uint32) {
	renderer.SetDrawColor(uint8(color >> 16), uint8(color >> 8), uint8(color), uint8(color >> 24))
	renderer.Clear()
}

// specify the stroke color
func Stroke(color uint32) {
	if color == 0 {
		color = 0xff000000
	}
	strokeColor = color
	strokeEnabled = true
	renderer.SetDrawColor(uint8(strokeColor >> 16), uint8(strokeColor >> 8), uint8(strokeColor), uint8(strokeColor >> 24))
}

// specify the fill color
func Fill(color uint32) {
	if color == 0 {
		color = 0xff000000
	}
	fillColor = color
	fillEnabled = true
	renderer.SetDrawColor(uint8(fillColor >> 16), uint8(fillColor >> 8), uint8(fillColor), uint8(fillColor >> 24))
}

// draw a point with specified position
func Point(x, y int) {
	if strokeEnabled {
		renderer.DrawPoint(x, y)
	}
}

// draw a rectangle with specified positions
func Line(x1, y1, x2, y2 int) {
	if strokeEnabled {
		renderer.DrawLine(x1, y1, x2, y2)
	}
}

// draw a rectangle with specified position and dimension
func Rect(x, y, w, h int) {
	rect := sdl.Rect {int32(x), int32(y), int32(w), int32(h)}
	if fillEnabled {
		renderer.FillRect(&rect)
	}
	if strokeEnabled {
		renderer.DrawRect(&rect)
	}
}

// load an image with specified filename (.jpg, .png, .gif, .bmp, .tiff, .webp, etc)
func LoadImage(filename string) PImage {
	surface := img.Load(filename)
	defer surface.Free()

	texture := renderer.CreateTextureFromSurface(surface)
	img := PImage {texture, int(surface.W), int(surface.H)}
	return img
}

// draw image normally on the screen
func Image(img PImage, x, y int) {
	rect := sdl.Rect {int32(x), int32(y), int32(img.Width), int32(img.Height)}
	renderer.Copy(img.texture, nil, &rect)
}

// draw image with specified size
func ImageScaled(img PImage, x, y, w, h int) {
	rect := sdl.Rect {int32(x), int32(y), int32(w), int32(h)}
	renderer.Copy(img.texture, nil, &rect)
}

// fill screen with the image scaled to the screen size
func ImageFill(img PImage) {
	renderer.Copy(img.texture, nil, nil)
}

package pogo

import (
	"github.com/jackyb/go-sdl2/sdl"
	"github.com/jackyb/go-sdl2/sdl_image"
	"github.com/jackyb/go-sdl2/sdl_mixer"
	"os"
	"log"
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
	surface *sdl.Surface
	texture *sdl.Texture
	Width, Height int
}

type Color sdl.Color

type PSound mix.Chunk

// internal variables
var window *sdl.Window = nil
var renderer *sdl.Renderer = nil
var strokeColor, fillColor Color = Color{0, 0, 0, 0xff}, Color{0xff, 0xff, 0xff, 0xff}
var strokeEnabled, fillEnabled = true, true
var frameRate = 60
var period = 1000 / frameRate

// public variables that users can use
var Width, Height int = 200, 200
var MouseX, MouseY int = 0, 0
var Key byte

// public assignable handler functions that users can use
var QuitHandler func() = nil
var KeyDownHandler, KeyUpHandler func() = nil, nil

func init() {
	window, renderer = sdl.CreateWindowAndRenderer(Width, Height, 0)
	if window == nil || renderer == nil {
		log.Fatal(sdl.GetError())
		os.Exit(1)
	}
	Background(200, 200, 200)

	endian := sdl.Endian()
	switch endian {
	case sdl.BIG_ENDIAN:
		if mix.OpenAudio(mix.DEFAULT_FREQUENCY, sdl.AUDIO_S16MSB,
				mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE) {
			log.Println(sdl.GetError())
		}
	case sdl.LIL_ENDIAN:
		if mix.OpenAudio(mix.DEFAULT_FREQUENCY, mix.DEFAULT_FORMAT,
				mix.DEFAULT_CHANNELS, mix.DEFAULT_CHUNKSIZE) {
			log.Println(sdl.GetError())
		}
	}
}

// setup the window and renderer
func Size(w, h int, flags uint32) {
	Width, Height = w, h

	if renderer == nil {
		os.Exit(1)
	}
	if window == nil {
		os.Exit(1)
	}

	renderer.Destroy()
	window.Destroy()

	window, renderer = sdl.CreateWindowAndRenderer(Width, Height, flags)
	Background(200, 200, 200)
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
func Background(r, g, b uint8) {
	renderer.SetDrawColor(r, g, b, 0xff)
	renderer.Clear()
}

func NoStroke() {
	strokeEnabled = false
}

// specify the stroke color
func Stroke(r, g, b, a uint8) {
	strokeColor = Color{r, g, b, a}
	strokeEnabled = true
}

func NoFill() {
	fillEnabled = false
}

// specify the fill color
func Fill(r, g, b, a uint8) {
	fillColor = Color {r, g, b, a}
	fillEnabled = true
}

// draw a point with specified position
func Point(x, y int) {
	renderer.SetDrawColor(strokeColor.R, strokeColor.G, strokeColor.B, strokeColor.A)
	if strokeEnabled {
		renderer.DrawPoint(x, y)
	}
}

// draw a rectangle with specified positions
func Line(x1, y1, x2, y2 int) {
	renderer.SetDrawColor(strokeColor.R, strokeColor.G, strokeColor.B, strokeColor.A)
	if strokeEnabled {
		renderer.DrawLine(x1, y1, x2, y2)
	}
}

// draw a rectangle with specified position and dimension
func Rect(x, y, w, h int) {
	rect := sdl.Rect {int32(x), int32(y), int32(w), int32(h)}
	if fillEnabled {
		renderer.SetDrawColor(fillColor.R, fillColor.G, fillColor.B, fillColor.A)
		renderer.FillRect(&rect)
	}
	if strokeEnabled {
		renderer.SetDrawColor(strokeColor.R, strokeColor.G, strokeColor.B, strokeColor.A)
		renderer.DrawRect(&rect)
	}
}

// load an image with specified filename (.jpg, .png, .gif, .bmp, .tiff, .webp, etc)
func LoadImage(filename string) PImage {
	var image PImage
	var surface *sdl.Surface
	var texture *sdl.Texture

	surface = img.Load(filename)
	if surface == nil {
		log.Println(sdl.GetError())
		goto out
	}
	surface = surface.ConvertFormat(sdl.PIXELFORMAT_ARGB8888, 0)
	texture = renderer.CreateTextureFromSurface(surface)
	if texture == nil {
		log.Println(sdl.GetError())
		goto out
	}
	image = PImage {surface, texture, int(surface.W), int(surface.H)}
out:
	return image
}

// draw image normally on the screen
func Image(image PImage, x, y int) {
	rect := sdl.Rect {int32(x), int32(y), int32(image.Width), int32(image.Height)}
	renderer.Copy(image.texture, nil, &rect)
}

// draw image with specified size
func ImageScaled(image PImage, x, y, w, h int) {
	rect := sdl.Rect {int32(x), int32(y), int32(w), int32(h)}
	renderer.Copy(image.texture, nil, &rect)
}

// fill screen with the image scaled to the screen size
func ImageFill(image PImage) {
	renderer.Copy(image.texture, nil, nil)
}

// get the image pixels in ARGB format
func (image *PImage) Pixels() []uint32 {
	return sdl.U8To32Array(image.surface.Pixels())
}

// update the image pixels for rendering
func (image *PImage) Update() {
	image.texture.Update(nil, image.surface.Data(), int(image.surface.W * int32(image.surface.Format.BytesPerPixel)))
}

func LoadSound(filename string) *PSound {
	return (*PSound) (mix.LoadWAV(filename))
}

func (sound *PSound) Play(loops int) bool {
	chunk := (*mix.Chunk) (sound)
	return chunk.PlayChannel(-1, loops)
}

func (sound *PSound) Playing() bool {
	return mix.Playing(-1)
}

func (sound *PSound) Free() {
	chunk := (*mix.Chunk) (sound)
	chunk.Free()
}

func Error() string {
	return sdl.GetError()
}

func PrintError() {
	log.Println(Error())
}

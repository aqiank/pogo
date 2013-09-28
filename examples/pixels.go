//package main

import p "github.com/jackyb/pogo"

var image p.PImage

func main() {
	p.Size(1024, 768, p.OPENGL)

	image = p.LoadImage("brush.png")
	pixels := image.Pixels()
	for i, _ := range pixels {
		pixels[i] |= (pixels[i] << 16)
	}
	image.Update()
	p.Loop(draw) // set the loop function
}

func draw() {
	p.Image(image, p.MouseX, p.MouseY)
}

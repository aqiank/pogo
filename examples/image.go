//package main // uncomment this line to try

import p "github.com/jackyb/pogo"

var image p.PImage

func main() {
	p.Size(1024, 768, p.OPENGL)

	image = p.LoadImage("brush.png")

	p.Loop(draw) // set the loop function
}

func draw() {
	p.Image(image, p.MouseX, p.MouseY)
}

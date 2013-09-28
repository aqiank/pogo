//package main

import p "github.com/jackyb/pogo"

var mode byte = '1'

func main() {
	p.Size(1024, 768, p.OPENGL)
	p.KeyUpHandler = func() {
		mode = p.Key
	}
	p.Fill(255, 0, 0, 255)
	p.Stroke(0, 0, 0, 255)

	p.Loop(draw) // set the loop function
}

func draw() {
	switch (mode) {
	case '1':
		p.Rect(p.MouseX, p.MouseY, 100, 100)
	case '2':
		p.Line(p.MouseX, p.MouseY, p.MouseX + 100, p.MouseY + 100)
	case '3':
		p.Point(p.MouseX, p.MouseY)
	}
}

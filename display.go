package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	running  bool
}

func (display *Display) init() int {
	var err error
	display.window, err = sdl.CreateWindow("GB", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		160, 144, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	display.renderer, err = sdl.CreateRenderer(display.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		return 2
	}
	display.running = true
	return 0
}

func (display *Display) display(gpu Gpu) int {

	display.running = true
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			display.running = false
		}
	}
	if gpu.rendering {
		for x := 1; x < 160; x++ {
			for y := 1; y < 144; y++ {
				if gpu.frame_buffer[x][y] > 0 {
					display.renderer.SetDrawColor(255, 255, 255, 255)
					display.renderer.DrawPoint(int32(x), int32(y))
				}
			}
		}
		display.renderer.Present()
	}
	return 0
}

func (display *Display) close() {
	display.renderer.Destroy()
	display.window.Destroy()
}

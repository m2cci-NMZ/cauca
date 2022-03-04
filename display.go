package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	window       *sdl.Window
	renderer     *sdl.Renderer
	vramWindow   *sdl.Window
	vramRenderer *sdl.Renderer
	running      bool
}

const tile_size int = 4

func (display *Display) init() int {
	var err error
	display.window, err = sdl.CreateWindow("GB", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		160*int32(tile_size), 144*int32(tile_size), sdl.WINDOW_SHOWN)
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

func (display *Display) initVramViewer() int {
	var err error
	display.vramWindow, err = sdl.CreateWindow("GB", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		8*16, 8*16, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %s\n", err)
		return 1
	}
	display.vramRenderer, err = sdl.CreateRenderer(display.vramWindow, -1, sdl.RENDERER_ACCELERATED)
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
		display.renderer.SetDrawColor(255, 255, 255, 0)
		display.renderer.Clear()
		for x := 1; x < 160*tile_size; x++ {
			for y := 1; y < 144*tile_size; y++ {
				if gpu.frame_buffer[x/tile_size][y/tile_size] > 0 {
					display.renderer.SetDrawColor(0, 0, 0, 255)
					display.renderer.DrawPoint(int32(x), int32(y))
				}
			}
		}
		display.renderer.Present()
	}
	return 0
}

func (display *Display) displayVram(gpu Gpu, mem Memory) int {

	display.running = true
	vram := gpu.getVram(mem)
	if gpu.rendering {
		display.vramRenderer.SetDrawColor(255, 255, 255, 0)
		display.vramRenderer.Clear()
		var screenx int = 0
		var screeny int = 0
		for y, slice := range vram {
			for x, pixel := range slice {
				if y%(8*16) == 0 && y != 0 {
					screenx += 8
					screeny += 8
				}
				if pixel > 0 {
					display.vramRenderer.SetDrawColor(0, 0, 0, 255)
					display.vramRenderer.DrawPoint(int32(screenx+x), int32(screeny+y%(8*16)))
				}
			}
		}
		display.vramRenderer.Present()
	}
	return 0
}

func (display *Display) close() {
	display.renderer.Destroy()
	display.window.Destroy()
}

func (display *Display) vramClose() {
	display.vramRenderer.Destroy()
	display.vramWindow.Destroy()
}

package main

import (
	"fmt"
	"os"
	//"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	emu "github.com/heyymrdj/chip8/emulator"
)

const BEEPSOUND = "./resources/beep.ogg"

func main() {
	modifier := 10
	if len(os.Args) < 2 {
		fmt.Println("Please provide a ROM path")
		os.Exit(1)
	}
	fileName := os.Args[1]
	fmt.Println(len(os.Args))

	c8 := emu.Init()
	if loadErr := c8.LoadProgram(fileName); loadErr != nil {
		panic(loadErr)
	}

	title := fmt.Sprint("CHIP8 Emulator - ", fileName)
	rl.InitWindow(640, 320, title)
	defer rl.CloseWindow()
	rl.InitAudioDevice() // Initialize audio device

	defer rl.CloseAudioDevice()      // De-initialize audio device
	sound := rl.LoadSound(BEEPSOUND) // Load audio file
	defer rl.UnloadSound(sound)      // Unload sound data

	// Function to be called for beep sound
	c8.AddBeep(func() {
		rl.PlaySound(sound)
	})

	// Cap FPS to 60 for propper speeds
	rl.SetTargetFPS(60)

	// Define colors
	foregroundColor := rl.NewColor(0, 0, 255, 255) // Blue
	backgroundColor := rl.NewColor(0, 0, 0, 255)   // Black

	// Game Loop
	for !rl.WindowShouldClose() {
		// Run CPU cycle
		c8.Cycle()
		// Draw screen if flag was set
		if c8.Draw() {
			rl.BeginDrawing()

			rl.ClearBackground(backgroundColor)

			// Get the display buffer and render
			vector := c8.Buffer()
			for j := 0; j < len(vector); j++ {
				for i := 0; i < len(vector[j]); i++ {
					color := backgroundColor
					if vector[j][i] != 0 {
						color = foregroundColor
					}
					rl.DrawRectangle(int32(i*modifier), int32(j*modifier), int32(modifier), int32(modifier), color)
				}
			}
			if c8.SoundTimer > 0 {
				c8.Beeper()
			}
			checkInput(&c8)
			rl.EndDrawing()
		}

	}
}

func checkInput(c8 *emu.Chip8) {
	// If KeyUp
	if rl.IsKeyUp(rl.KeyOne) {
		c8.Key(0x1, false)
	}
	if rl.IsKeyUp(rl.KeyTwo) {
		c8.Key(0x2, false)
	}
	if rl.IsKeyUp(rl.KeyThree) {
		c8.Key(0x3, false)
	}
	if rl.IsKeyUp(rl.KeyFour) {
		c8.Key(0xC, false)
	}
	if rl.IsKeyUp(rl.KeyQ) {
		c8.Key(0x4, false)
	}
	if rl.IsKeyUp(rl.KeyW) {
		c8.Key(0x5, false)
	}
	if rl.IsKeyUp(rl.KeyE) {
		c8.Key(0x6, false)
	}
	if rl.IsKeyUp(rl.KeyR) {
		c8.Key(0xD, false)
	}
	if rl.IsKeyUp(rl.KeyA) {
		c8.Key(0x7, false)
	}
	if rl.IsKeyUp(rl.KeyS) {
		c8.Key(0x8, false)
	}
	if rl.IsKeyUp(rl.KeyD) {
		c8.Key(0x9, false)
	}
	if rl.IsKeyUp(rl.KeyF) {
		c8.Key(0xE, false)
	}
	if rl.IsKeyUp(rl.KeyZ) {
		c8.Key(0xA, false)
	}
	if rl.IsKeyUp(rl.KeyX) {
		c8.Key(0x0, false)
	}
	if rl.IsKeyUp(rl.KeyC) {
		c8.Key(0xB, false)
	}
	if rl.IsKeyUp(rl.KeyV) {
		c8.Key(0xF, false)
	}

	// If KeyDown
	if rl.IsKeyDown(rl.KeyOne) {
		c8.Key(0x1, true)
	}
	if rl.IsKeyDown(rl.KeyTwo) {
		c8.Key(0x2, true)
	}
	if rl.IsKeyDown(rl.KeyThree) {
		c8.Key(0x3, true)
	}
	if rl.IsKeyDown(rl.KeyFour) {
		c8.Key(0xC, true)
	}
	if rl.IsKeyDown(rl.KeyQ) {
		c8.Key(0x4, true)
	}
	if rl.IsKeyDown(rl.KeyW) {
		c8.Key(0x5, true)
	}
	if rl.IsKeyDown(rl.KeyE) {
		c8.Key(0x6, true)
	}
	if rl.IsKeyDown(rl.KeyR) {
		c8.Key(0xD, true)
	}
	if rl.IsKeyDown(rl.KeyA) {
		c8.Key(0x7, true)
	}
	if rl.IsKeyDown(rl.KeyS) {
		c8.Key(0x8, true)
	}
	if rl.IsKeyDown(rl.KeyD) {
		c8.Key(0x9, true)
	}
	if rl.IsKeyDown(rl.KeyF) {
		c8.Key(0xE, true)
	}
	if rl.IsKeyDown(rl.KeyZ) {
		c8.Key(0xA, true)
	}
	if rl.IsKeyDown(rl.KeyX) {
		c8.Key(0x0, true)
	}
	if rl.IsKeyDown(rl.KeyC) {
		c8.Key(0xB, true)
	}
	if rl.IsKeyDown(rl.KeyV) {
		c8.Key(0xF, true)
	}
}

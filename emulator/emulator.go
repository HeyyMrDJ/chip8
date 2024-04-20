package emulator

import (
	"fmt"
	"math/rand"
	"os"
)

var fontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

type Chip8 struct {
	display        [32][64]uint8
	memory         [4096]uint8
	registers      [16]uint8
	keys           [16]uint8
	stack          [16]uint16
	opCode         uint16
	sp             uint16
	programCounter uint16
	index          uint16
	delayTimer     uint8
	SoundTimer     uint8
	shouldDraw     bool
	Beeper         func()
}

func Init() Chip8 {
	chip8 := Chip8{
		shouldDraw:     true,
		programCounter: 0x200,
		Beeper:         func() {},
	}

	//for i := 0; i < len(fontSet); i++ {
	//	chip8.memory[i] = fontSet[i]
	//}

	return chip8
}

func (c *Chip8) Buffer() [32][64]uint8 {

	return c.display
}

func (c *Chip8) Draw() bool {
	sd := c.shouldDraw
	c.shouldDraw = false

	return sd
}

func (c *Chip8) AddBeep(fn func()) {
	c.Beeper = fn
}

func (c *Chip8) Key(num uint8, down bool) {
	if down {
		c.keys[num] = 1
	} else {
		c.keys[num] = 0
	}
}

func (c *Chip8) LoadProgram(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fStat, err := file.Stat()
	if err != nil {
		return err
	}

	if int64(len(c.memory)-512) < fStat.Size() {
		return fmt.Errorf("Program too large for memory")
	}

	buffer := make([]byte, fStat.Size())
	if _, err := file.Read(buffer); err != nil {
		return err
	}

	for i := 0; i < len(buffer); i++ {
		c.memory[i+512] = buffer[i]
	}

	return nil
}

func (c *Chip8) Cycle() {
	fmt.Println(c.opCode)
	if c.programCounter > 0xFFF {
		panic("Invalid opcode")
	}

	c.opCode = uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])

	switch c.opCode >> 12 {
	case 0x0:
		switch c.opCode & 0x000F {
		case 0x0:
			fmt.Println("New 0x0")
			for y := range c.display {
				for x := range c.display[y] {
					c.display[y][x] = uint8(0)
				}
			}
			c.shouldDraw = true
		case 0x000E:
			fmt.Println("New 0x000E")
			c.sp -= 1
			c.programCounter = c.stack[c.sp]
		default:
			fmt.Printf("Unknown opcode %X\n", c.opCode)
		}
		c.programCounter += 2
	case 0x1:
		c.programCounter = c.opCode & 0x0FFF
	case 0x2:
		c.stack[c.sp] = c.programCounter
		c.sp += 1
		c.programCounter = c.opCode & 0x0FFF
	case 0x3:
		x := c.opCode & 0x0F00 >> 8
		if c.registers[x] == uint8(c.opCode&0x00FF) {
			c.programCounter += 2
		}

		c.programCounter += 2
	case 0x4:
		x := (c.opCode & 0x0F00) >> 8
		if c.registers[x] != uint8((c.opCode & 0x00FF)) {
			c.programCounter += 2
		}

		c.programCounter += 2
	case 0x5:
		x := (c.opCode & 0x0F00) >> 8
		y := (c.opCode & 0x00F0) >> 4

		if c.registers[x] == c.registers[y] {
			c.programCounter += 2
		}
		c.programCounter += 2
	case 0x6:
		x := (c.opCode & 0x0F00) >> 8
		c.registers[x] = uint8(c.opCode & 0x00FF)
		c.programCounter += 2
	case 0x7:
		x := (c.opCode & 0x0F00) >> 8
		c.registers[x] += uint8(c.opCode & 0x00FF)
		c.programCounter += 2
	case 0x8:
		x := c.opCode & 0x0F00 >> 8
		y := c.opCode & 0x00F0 >> 4
		z := c.opCode & 0x000F

		switch z {
		case 0:
			c.registers[x] = c.registers[y]
		case 1:
			c.registers[x] |= c.registers[y]
		case 2:
			c.registers[x] &= c.registers[y]
		case 3:
			c.registers[x] ^= c.registers[y]
		case 4:
			c.registers[x] += c.registers[y]
			sum := uint16(c.registers[x])
			sum += uint16(c.registers[y])

			if c.registers[y] > (0xFF - c.registers[x]) {
				c.registers[0xF] = 1
			} else {
				c.registers[0xf] = 0
			}
		case 5:
			if c.registers[y] > c.registers[x] {
				c.registers[0xF] = 0
			} else {
				c.registers[0xF] = 1
			}
			c.registers[x] -= c.registers[y]
		case 6:
			c.registers[0xF] = c.registers[x] & 0x1
			c.registers[x] >>= 1
		case 7:
			if c.registers[x] > c.registers[y] {
				c.registers[0xF] = 0
			} else {
				c.registers[0xF] = 1
			}
			c.registers[x] = c.registers[y] - c.registers[x]
		case 0xE:
			c.registers[0xF] = c.registers[x] >> 7
			c.registers[x] <<= 1

		default:
			fmt.Printf("Current ALU opcode: %x", c.opCode)

		}
		c.programCounter += 2
	case 0x9:
		x := c.opCode & 0x0F00 >> 8
		y := c.opCode & 0x00F0 >> 4
		if c.registers[x] != c.registers[y] {
			c.programCounter += 2
		}
		c.programCounter += 2
	case 0xA:
		c.index = c.opCode & 0x0FFF
		c.programCounter += 2
	case 0xB:
		c.programCounter = (c.opCode & 0x0FFF) + uint16(c.registers[0])
	case 0xC:
		x := c.opCode & 0x0F00 >> 8
		kk := c.opCode & 0x00FF

		c.registers[x] = uint8(rand.Intn(256)) & uint8(c.opCode&kk)
		c.programCounter += 2
	case 0xD:
		c.registers[0xF] = 0
		regX := c.registers[(c.opCode&0x0F00)>>8]
		regY := uint16(c.registers[(c.opCode&0x00F0)>>4])
		height := c.opCode & 0x000F
		var j uint16 = 0
		var i uint16 = 0
		for j = 0; j < height; j++ {
			if regY+j >= 32 {
				continue
			}
			pixel := c.memory[c.index+j]
			for i = 0; i < 8; i++ {
				if uint16(regX)+i >= 64 {
					continue
				}
				if (pixel & (0x80 >> i)) != 0 {
					if c.display[(uint8(regY) + uint8(j))][regX+uint8(i)] == 1 {
						c.registers[0xF] = 1
					}
					c.display[(uint8(regY) + uint8(j))][regX+uint8(i)] ^= 1
				}
			}
		}
		c.shouldDraw = true
		c.programCounter = c.programCounter + 2
	case 0xE:
		x := c.opCode & 0x0F00 >> 8
		m := c.opCode & 0x00FF

		if m == 0x9E {
			if c.keys[c.registers[x]] != 0 {
				c.programCounter += 2
			}
		} else if m == 0xA1 {
			if c.keys[c.registers[x]] == 0 {
				c.programCounter += 2
			}
		}
		c.programCounter += 2
	case 0xF:
		x := c.opCode & 0x0F00 >> 8
		m := c.opCode & 0x00FF

		if m == 0x07 {
			c.registers[x] = c.delayTimer
		} else if m == 0x0A {
			keyPressed := false

			i := uint8(0)
			for i = 0; i < 16; i++ {
				if c.keys[i] != 0 {
					c.registers[x] = uint8(i)
					keyPressed = true
				}
			}
			if !keyPressed {
				return
			}
		} else if m == 0x15 {
			c.delayTimer = c.registers[x]
		} else if m == 0x18 {
			c.SoundTimer = c.registers[x]
		} else if m == 0x1E {
			if c.index+uint16(c.registers[x]) > 0xFFF {
				c.registers[0xF] = 1
			} else {
				c.registers[0xF] = 0
			}
			c.index += uint16(c.registers[x])
		} else if m == 0x29 {
			c.index = uint16(c.registers[x] * 0x5)
		} else if m == 0x33 {
			c.memory[c.index] = c.registers[x] / 100
			c.memory[c.index+1] = (c.registers[x] / 10) % 10
			c.memory[c.index+2] = c.registers[x] % 10
		} else if m == 0x55 {
			i := uint16(0)
			for i = 0; i <= x; i++ {
				c.memory[c.index+i] = c.registers[i]
			}
			c.index += uint16(c.registers[x]) + 1
		} else if m == 0x65 {
			i := uint16(0)
			for i = 0; i <= x; i++ {
				c.registers[i] = c.memory[c.index+i]
			}
			c.index += uint16(c.registers[x]) + 1

		}
		c.programCounter += 2

	default:
		fmt.Printf("Current opcode: %X", c.opCode)
	}

	if c.delayTimer > 0 {
		c.delayTimer -= 1
	}
	if c.SoundTimer > 0 {
		c.SoundTimer -= 1
		if c.SoundTimer == 1 {
			c.Beeper()
		}
	}
}

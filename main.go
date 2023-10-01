package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 64
	screenHeight = 32
)

type Chip8 struct {
	memory     [4096]byte
	v          [16]byte
	i          uint16
	pc         uint16
	stack      [16]uint16
	sp         int
	display    [screenWidth * screenHeight]byte
	key        [16]bool
	delay      uint8
	sound      uint8
	shouldDraw bool
	opcode     uint16
}

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

func (c *Chip8) LoadROM(romPath string) error {
	romData, err := ebitenutil.OpenFile(romPath)
	if err != nil {
		return err
	}
	defer romData.Close()

	n, err := romData.Read(c.memory[0x200:])
	if err != nil {
		return err
	}

	fmt.Printf("Loaded %d bytes from the ROM file.\n", n)

	return nil
}

func (c *Chip8) emulateOpcode(opcode uint16) {
	x := int((opcode & 0x0F00) >> 8)
	y := int((opcode & 0x00F0) >> 4)
	kk := byte(opcode & 0x00FF)
	nnn := opcode & 0x0FFF

	switch opcode & 0xF000 {
	case 0x00E0:
		// Efface l'écran
		for i := range c.display {
			c.display[i] = 0
		}
		c.shouldDraw = true
		c.pc += 2

	case 0x00EE:
		// Retour de sous-routine
		if c.sp > 0 {
			c.sp--
			c.pc = c.stack[c.sp]
		}

	default:
		fmt.Printf("Opcode non pris en charge: %X\n", opcode)
	}

	switch opcode & 0xF000 {
	case 0x1000:
		// Saute à l'adresse NNN
		c.pc = nnn

	case 0x2000:
		// Appel de sous-routine
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = nnn

	case 0x3000:
		// Sauter si VX est égal à KK
		if c.v[x] == kk {
			c.pc += 2
		}

	case 0x4000:
		// Sauter si VX n'est pas égal à KK
		if c.v[x] != kk {
			c.pc += 2
		}

	case 0x5000:
		// Sauter si VX est égal à VY
		if c.v[x] == c.v[y] {
			c.pc += 2
		}

	case 0x6000:
		// Mettre KK dans VX
		c.v[x] = kk

	case 0x7000:
		// Ajouter KK à VX
		c.v[x] += kk

	case 0x8000:
		switch opcode & 0x000F {
		case 0x0000:
			// Mettre VY dans VX
			c.v[x] = c.v[y]

		case 0x0001:
			// VX = VX | VY
			c.v[x] |= c.v[y]

		case 0x0002:
			// VX = VX & VY
			c.v[x] &= c.v[y]

		case 0x0003:
			// VX = VX ^ VY
			c.v[x] ^= c.v[y]

		case 0x0004:
			// Addition avec retenue
			if int(c.v[x])+int(c.v[y]) > 255 {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[x] += c.v[y]

		case 0x0005:
			// Soustraction avec retenue
			if c.v[x] > c.v[y] {
				c.v[0xF] = 1
			} else {
				c.v[0xF] = 0
			}
			c.v[x] -= c.v[y]

		case 0x0006:
			// Décalage à droite
			c.v[0xF] = c.v[x] & 0x01
			c.v[x] >>= 1

		case 0x0007:
			// VX = VY - VX
			if c.v[y] > c.v[x] {
				c.v[0xF] = 0
			} else {
				c.v[0xF] = 1
			}
			c.v[x] = c.v[y] - c.v[x]

		case 0x000E:
			// Décalage à gauche
			c.v[0xF] = (c.v[x] >> 7) & 0x01
			c.v[x] <<= 1

		default:
			fmt.Printf("Opcode non pris en charge: %X\n", opcode)
		}

	case 0x9000:
		// Sauter si VX n'est pas égal à VY
		if c.v[x] != c.v[y] {
			c.pc += 2
		}

	case 0xA000:
		// Mettre NNN dans I
		c.i = nnn

	case 0xB000:
		// Sauter à l'adresse NNN + V0
		c.pc = nnn + uint16(c.v[0])

	case 0xC000:
		// Mettre un nombre aléatoire ET KK dans VX
		randomValue := byte(rand.Intn(256))
		c.v[x] = randomValue & kk

	case 0xD000:
		// Dessine un sprite à l'écran
		xPos := int(c.v[x])
		yPos := int(c.v[y])
		height := int(opcode & 0x000F)
		collision := false

		for row := 0; row < height; row++ {
			pixel := c.memory[c.i+uint16(row)]
			y := (yPos + row) % screenHeight

			for col := 0; col < 8; col++ {
				x := (xPos + col) % screenWidth

				if (pixel & (0x80 >> col)) != 0 {
					if c.display[y*screenWidth+x] == 1 {
						collision = true
					}
					c.display[y*screenWidth+x] ^= 1
				}
			}
		}

		if collision {
			c.v[0xF] = 1
		} else {
			c.v[0xF] = 0
		}

		c.shouldDraw = true
		c.pc += 2

	default:
		fmt.Printf("Opcode non pris en charge: %X\n", opcode)
	}
}

func (c *Chip8) Draw(screen *ebiten.Image) {
	if c.shouldDraw {
		for y := 0; y < screenHeight; y++ {
			for x := 0; x < screenWidth; x++ {
				// Dessinez un pixel blanc si c.display[y*screenWidth+x] vaut 1, sinon noir
				if c.display[y*screenWidth+x] == 1 {
					screen.Set(x, y, color.White)
				} else {
					screen.Set(x, y, color.Black)
				}
			}
		}
		c.shouldDraw = false
	}
}

func (c *Chip8) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (c *Chip8) Update() error {
	// Lire l'opcode à partir de la mémoire à la position actuelle du PC
	opcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])

	// Appeler emulateOpcode pour émuler l'opcode
	c.emulateOpcode(opcode)

	// La logique de mise à jour de l'émulateur va ici

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(screenWidth*10, screenHeight*10)
	ebiten.SetWindowTitle("CHIP-8 Emulator")

	chip8 := &Chip8{}

	// Chargez la ROM que vous souhaitez exécuter
	romPath := "roms8/filter.ch8"
	if err := chip8.LoadROM(romPath); err != nil {
		log.Fatal(err)
	}

	// Démarrez Ebiten
	if err := ebiten.RunGame(chip8); err != nil {
		log.Fatal(err)
	}
}

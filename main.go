package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

type Chip8 struct {
	memory  [4096]byte
	v       [16]byte
	i       uint16
	pc      uint16
	stack   []uint16
	sp      int
	display [64 * 32]byte
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowTitle("Emu8")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 64
	screenHeight = 32
)

type GoChip8 struct {
	memory  [4096]byte
	v       [16]byte
	i       uint16
	pc      uint16
	stack   []uint16
	sp      int
	display [screenWidth * screenHeight]byte
	key     [16]bool
	delay   uint8
	sound   uint8
}

func NewGoChip8() *GoChip8 {
	return &GoChip8{}
}

func (c8 *GoChip8) LoadROM(filename string) error {
	file, err := os.Open("roms8/c8_test.c8")
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func main() {

	goChip8 := NewGoChip8()

	if len(os.Args) != 2 {
		fmt.Println("roms8/c8_test.c8")
		return
	}
	if err := goChip8.LoadROM(os.Args[1]); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth*10, screenHeight*10)
	ebiten.SetWindowTitle("Chip-8 Emulator")
	if err := ebiten.RunGame(goChip8); err != nil {
		log.Fatal(err)
	}
}

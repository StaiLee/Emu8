package main
// Importation des packages
import (
	"emu8/emul8"
	"log"
	"os"
	"github.com/hajimehoshi/ebiten/v2"
)
// Création fonction main
func main() {
	// verifier les arg
	if len(os.Args) < 2 {   
		// erreur case 
		log.Fatal("roms8/file.ch8") // 
	}
	programName := os.Args[1] 
	// création de'object chip
	chip := emul8.InitiateChip8()	
	chip.LoadGUI(programName)
	// set la size
	ebiten.SetWindowSize(64*10, 32*10)
	// titrer la fenetre
	ebiten.SetWindowTitle(programName)
	// check error
	if err := ebiten.RunGame(&chip); err != nil {
		log.Fatal(err)
	}
}

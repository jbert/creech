package main

import (
	"log"

	"github.com/jbert/creech"
	"github.com/jbert/creech/render"
)

func main() {
	r := render.NewScreen()
	game := creech.NewGame(r)
	err := game.Init()
	if err != nil {
		log.Fatalf("Init with error: %s", err)
	}
	err = game.Run()
	if err != nil {
		log.Fatalf("Exit with error: %s", err)
	}
}

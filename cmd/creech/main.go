package main

import (
	"flag"
	"log"

	"github.com/jbert/creech"
	"github.com/jbert/creech/render"
)

type options struct {
	renderMode string
	hostPort   string
}

func flagsToOptions() *options {
	var o options
	flag.StringVar(&o.renderMode, "render", "screen", "render mode: 'screen' or 'web'")
	flag.StringVar(&o.hostPort, "hostport", ":8080", "host:port for web mode")
	flag.Parse()
	return &o
}

func main() {
	o := flagsToOptions()

	var r render.Renderer
	switch o.renderMode {
	case "screen":
		r = render.NewScreen()
	case "web":
		r = render.NewWeb(o.hostPort)
	default:
		log.Fatalf("Unknown render mode: %s", o.renderMode)
	}

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

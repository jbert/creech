package main

import (
	"flag"
	"log"
	"time"

	"github.com/jbert/creech"
	"github.com/jbert/creech/render"
)

type options struct {
	renderMode string
	hostPort   string
	tick       time.Duration
}

func flagsToOptions() *options {
	var o options
	flag.StringVar(&o.renderMode, "render", "screen", "render mode: 'screen' or 'web'")
	flag.StringVar(&o.hostPort, "hostport", ":8080", "host:port for web mode")
	flag.DurationVar(&o.tick, "tick", time.Second, "Tick duration")
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

	game := creech.NewGame(r, o.tick)
	err := game.Init()
	if err != nil {
		log.Fatalf("Init with error: %s", err)
	}
	err = game.Run()
	if err != nil {
		log.Fatalf("Exit with error: %s", err)
	}
}

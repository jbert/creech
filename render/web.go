package render

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/jbert/creech/pos"
)

type Web struct {
	hostport       string
	mux            *http.ServeMux
	width          float64
	height         float64
	pixelsPerMetre float64

	rootTemplate *template.Template
	drawCh       chan drawCommand
}

type drawType int

const (
	startFrame drawType = iota
	drawPoly
	finishFrame
)

type drawCommand struct {
	What   drawType
	Points []pos.Pos
}

//go:embed static/root.html
var rootTemplateString string

func NewWeb(hostport string) *Web {
	return &Web{
		hostport: hostport,
		mux:      http.NewServeMux(),

		pixelsPerMetre: 20.0,

		rootTemplate: template.Must(template.New("root").Parse(rootTemplateString)),
		drawCh:       make(chan drawCommand),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (w *Web) handleWebSocket(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Can't upgrade websocket: %s", err), http.StatusInternalServerError)
		return
	}
	for cmd := range w.drawCh {
		wsWriter, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Can't get websocket writer: %s", err), http.StatusInternalServerError)
			return
		}
		jsonCmd, err := json.Marshal(cmd)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Can't marshal to JSON: %s", err), http.StatusInternalServerError)
			err = wsWriter.Close()
			return
		}
		wsWriter.Write(jsonCmd)
		err = wsWriter.Close()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Can't close websocket writer: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

func (w *Web) handleRoot(rw http.ResponseWriter, r *http.Request) {
	tmplData := struct {
		WidthPixels  int
		HeightPixels int
		Scale        float64
		WSURL        string
		StartFrame   int
		FinishFrame  int
		DrawPoly     int
	}{
		int(w.pixelsPerMetre * w.width),
		int(w.pixelsPerMetre * w.height),
		w.pixelsPerMetre,
		"ws",
		int(startFrame),
		int(finishFrame),
		int(drawPoly),
	}
	err := w.rootTemplate.Execute(rw, tmplData)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Can't render template: %s", err), http.StatusInternalServerError)
		return
	}
	w.rootTemplate.Execute(os.Stdout, tmplData)
}

func (w *Web) Init(width, height float64) error {
	w.width = width
	w.height = height

	w.mux.HandleFunc("/", w.handleRoot)
	w.mux.HandleFunc("/ws", w.handleWebSocket)
	go func() {
		log.Printf("About to listen on [%s]", w.hostport)
		err := http.ListenAndServe(w.hostport, w.mux)
		if err != nil {
			log.Printf("Failed to ListenAndServe: %s\n", err)
		}
	}()
	return nil
}

func (w *Web) StartFrame() error {
	select {
	// Just drop it if we have no connection
	case w.drawCh <- drawCommand{What: startFrame}:
		return nil
	}
}

func (w *Web) FinishFrame() error {
	select {
	// Just drop it if we have no connection
	case w.drawCh <- drawCommand{What: finishFrame}:
		return nil
	}
}

func (w *Web) Draw(d Drawable) error {
	points := d.Web()
	select {
	// Just drop it if we have no connection
	case w.drawCh <- drawCommand{What: drawPoly, Points: points}:
		return nil
	}
}

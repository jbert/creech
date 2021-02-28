package render

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Web struct {
	hostport       string
	mux            *http.ServeMux
	width          float64
	height         float64
	pixelsPerMetre float64

	rootTemplate *template.Template
	drawCh       chan DrawCommand
}

//go:embed static/root.html
var rootTemplateString string

func NewWeb(hostport string) *Web {
	return &Web{
		hostport: hostport,
		mux:      http.NewServeMux(),

		pixelsPerMetre: 20.0,

		rootTemplate: template.Must(template.New("root").Parse(rootTemplateString)),
		drawCh:       make(chan DrawCommand),
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
		err = conn.WriteJSON(cmd)
		if err != nil {
			log.Printf("Can't writeJSON to websocket: %s", err)
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
		int(StartFrame),
		int(FinishFrame),
		int(DrawPoly),
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
	case w.drawCh <- DrawCommand{What: StartFrame}:
	}
	return nil
}

func (w *Web) FinishFrame() error {
	select {
	// Just drop it if we have no connection
	case w.drawCh <- DrawCommand{What: FinishFrame}:
	}
	return nil
}

func (w *Web) Draw(d Drawable) error {
	cmds := d.Web()
	for _, cmd := range cmds {
		select {
		// Just drop it if we have no connection
		case w.drawCh <- cmd:
		}
	}
	return nil
}

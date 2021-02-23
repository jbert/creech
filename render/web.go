package render

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type Web struct {
	hostport       string
	mux            *http.ServeMux
	width          float64
	height         float64
	pixelsPerMetre float64

	rootTemplate *template.Template
}

//go:embed static/root.html
var rootTemplateString string

func NewWeb(hostport string) *Web {
	return &Web{
		hostport: hostport,
		mux:      http.NewServeMux(),

		pixelsPerMetre: 5.0,

		rootTemplate: template.Must(template.New("root").Parse(rootTemplateString)),
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
	tickDur := time.Second
	ticker := time.NewTicker(tickDur)
	defer ticker.Stop()
	for {
		log.Printf("Sending ws message")
		wsWriter, err := conn.NextWriter(websocket.TextMessage)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Can't get websocket writer: %s", err), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(wsWriter, "Hi - it is %s\n", time.Now())
		err = wsWriter.Close()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Can't close websocket writer: %s", err), http.StatusInternalServerError)
			return
		}
		log.Printf("Sent ws message")
		<-ticker.C
	}
}

func (w *Web) handleRoot(rw http.ResponseWriter, r *http.Request) {
	tmplData := struct {
		CanvasID     string
		WidthPixels  int
		HeightPixels int
		WSURL        string
	}{
		"world",
		int(w.pixelsPerMetre * w.width),
		int(w.pixelsPerMetre * w.height),
		"ws",
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
	return nil
}

func (w *Web) FinishFrame() error {
	return nil
}

func (w *Web) DrawAt(x, y float64, d Drawable) error {
	return nil
}

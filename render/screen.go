package render

import (
	"io"
	"os"
)

type Screen struct {
	width, height int
	buffer        [][]byte
}

func NewScreen() *Screen {
	return &Screen{}
}

func (s *Screen) Init(w, h int) error {
	s.width = w
	s.height = h

	return nil
}

func (s *Screen) resetBuffer() {
	s.buffer = make([][]byte, s.height)
	for j := range s.buffer {
		s.buffer[j] = make([]byte, s.width+1)
		for i := range s.buffer[j] {
			s.buffer[j][i] = 0x20 // space
		}
		s.buffer[j][s.width] = 0x0a // \n
	}
	return
}

func (s *Screen) StartFrame() error {
	s.resetBuffer()
	return nil
}

func (s *Screen) FinishFrame() error {
	s.clearScreen(os.Stdout)
	s.writeBuffer(os.Stdout)
	return nil
}

func (s *Screen) DrawAt(i, j int, d Drawable) error {
	b := d.Screen()
	j = (-j + 3*s.height/2) % s.height
	i = (i + 3*s.width/2) % s.width
	s.buffer[j][i] = b
	return nil
}

func (s *Screen) clearScreen(w io.Writer) {
	// From 'clear | hd'
	clearByteSeq := []byte{0x1b, 0x5b, 0x48, 0x1b, 0x5b, 0x32, 0x4a, 0x1b, 0x5b, 0x33, 0x4a}
	w.Write(clearByteSeq)
}

func (s *Screen) writeBuffer(w io.Writer) {
	for i := range s.buffer {
		w.Write(s.buffer[i])
	}
}

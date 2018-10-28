package view

import (
	"fmt"
	"sync"
	"time"

	"github.com/aestek/baleno/buffer"
	"github.com/aestek/baleno/keymap"
	"github.com/aestek/baleno/state"
	log "github.com/sirupsen/logrus"
)

type StatusView struct {
	drawBuffer     DrawBuffer
	drawBufferLock sync.Mutex
	buffer         buffer.Buffer
	width, height  int
	state          *state.State
	infoState      *state.State

	infoCursorX int
	infoCursorY int
}

func NewStatusView(infoState *state.State) *StatusView {
	v := &StatusView{
		infoState: infoState,
	}

	go func() {
		for {
			c := make(chan state.Event, 10)
			err := infoState.Watch("window.panes.focused.view.cursor", c)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			for e := range c {
				log.Info(e)
				v.infoCursorX = infoState.MustGet("window.panes.focused.view.cursor.x").(int)
				v.infoCursorY = infoState.MustGet("window.panes.focused.view.cursor.y").(int)
				v.draw()
			}
		}
	}()

	return v
}

func (v *StatusView) Buffer() DrawBuffer {
	v.drawBufferLock.Lock()
	defer v.drawBufferLock.Unlock()
	return v.drawBuffer
}

func (v *StatusView) SetSize(w, h int) {
	if w == v.width && h == v.height {
		return
	}
	v.drawBufferLock.Lock()
	defer v.draw()
	defer v.drawBufferLock.Unlock()
	v.width = w
	v.height = h
	v.drawBuffer = NewDrawBuffer(w, h)
}

func (v *StatusView) Attach(s *state.State) {
	v.state = s
}

func (v *StatusView) HandleKeyPress(k keymap.KeyPress) {
}

func (v *StatusView) draw() {
	line := fmt.Sprintf("Ln %d, Col %d", v.infoCursorY+1, v.infoCursorX+1)

	v.drawBuffer.Clear()
	for i, r := range []rune(line) {
		v.drawBuffer[0][i].Char = r
	}
}

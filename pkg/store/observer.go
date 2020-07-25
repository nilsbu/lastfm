package store

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type observer interface {
	RequestRead(r rsrc.Locator)
	NotifyRead(r rsrc.Locator)
	RequestWrite(r rsrc.Locator)
	NotifyWrite(r rsrc.Locator)
	RequestRemove(r rsrc.Locator)
	NotifyRemove(r rsrc.Locator)
}

type counter struct {
	fChan chan<- format.Formatter
	eChan chan obEvent
	back  chan bool
	count map[rune][2]int
}

type obEvent struct {
	kind    rune
	receive bool
}

func NewObserver(fChan chan<- format.Formatter) observer {
	o := &counter{
		fChan: fChan,
		count: make(map[rune][2]int),
		eChan: make(chan obEvent),
		back:  make(chan bool),
	}

	go func(o *counter) {
		for e := range o.eChan {
			if e.receive {
				o.count[e.kind] = [2]int{o.count[e.kind][0] + 1, o.count[e.kind][1]}
			} else {
				o.count[e.kind] = [2]int{o.count[e.kind][0], o.count[e.kind][1] + 1}
			}
			o.sendFormat()
			o.back <- true
		}
	}(o)

	return o
}

func (o *counter) RequestRead(r rsrc.Locator) {
	o.eChan <- obEvent{'r', false}
	<-o.back
}

func (o *counter) NotifyRead(r rsrc.Locator) {
	o.eChan <- obEvent{'r', true}
	<-o.back
}

func (o *counter) RequestWrite(r rsrc.Locator) {
	o.eChan <- obEvent{'w', false}
	<-o.back
}

func (o *counter) NotifyWrite(r rsrc.Locator) {
	o.eChan <- obEvent{'w', true}
	<-o.back
}

func (o *counter) RequestRemove(r rsrc.Locator) {
	o.eChan <- obEvent{'d', false}
	<-o.back
}

func (o *counter) NotifyRemove(r rsrc.Locator) {
	o.eChan <- obEvent{'d', true}
	<-o.back
}

func (o *counter) sendFormat() {
	msg := fmt.Sprintf("r: %v/%v, w: %v/%v, rm: %v/%v",
		o.count['r'][0], o.count['r'][1],
		o.count['w'][0], o.count['w'][1],
		o.count['d'][0], o.count['d'][1],
	)

	o.fChan <- &format.Message{
		Msg: msg,
	}
}

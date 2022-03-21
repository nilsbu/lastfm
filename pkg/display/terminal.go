package display

import (
	"io"
	"os"
	"time"

	"github.com/nilsbu/lastfm/pkg/format"
)

type terminal struct {
	writer io.Writer
}

func NewTerminal() Display {
	return &terminal{writer: os.Stdout}
}

func (d *terminal) Display(f format.Formatter) error {
	return f.Plain(d.writer)
}

type timedTerminal struct {
	terminal
	timedF <-chan format.Formatter
	fCache format.Formatter
	fChan  chan format.Formatter
	eChan  chan error
}

func NewTimedTerminal(
	timedF <-chan format.Formatter,
	period time.Duration,
) Display {
	d := &timedTerminal{
		terminal: terminal{writer: os.Stdout},
		timedF:   timedF,
		fChan:    make(chan format.Formatter),
		eChan:    make(chan error),
	}

	go func(period time.Duration) {
		lastT := time.Now()
		for {
			select {
			case f := <-d.timedF:
				d.fCache = f

				now := time.Now()
				if now.Sub(lastT) >= period {
					if d.fCache != nil {
						if d.terminal.Display(d.fCache) != nil {
							return
						}
						d.fCache = nil
					}
					lastT = now
				}
			case f := <-d.fChan:
				err := d.terminal.Display(f)
				d.eChan <- err

			}
		}

	}(period)

	return d
}

func (d *timedTerminal) Display(f format.Formatter) error {
	d.fChan <- f
	return <-d.eChan
}

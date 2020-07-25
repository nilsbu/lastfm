package display

import (
	"io"
	"os"
	"time"

	"github.com/nilsbu/lastfm/pkg/format"
)

type Terminal struct {
	Writer io.Writer
}

func NewTerminal() *Terminal {
	return &Terminal{Writer: os.Stdout}
}

func (d *Terminal) Display(f format.Formatter) error {
	return f.Plain(d.Writer)
}

type TimedTerminal struct {
	Terminal
	timedF <-chan format.Formatter
	fCache format.Formatter
	fChan  chan format.Formatter
	eChan  chan error
}

func NewTimedTerminal(
	timedF <-chan format.Formatter,
	period time.Duration,
) *TimedTerminal {
	d := &TimedTerminal{
		Terminal: Terminal{Writer: os.Stdout},
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
						if d.Terminal.Display(d.fCache) != nil {
							return
						}
						d.fCache = nil
					}
					lastT = now
				}
			case f := <-d.fChan:
				err := d.Terminal.Display(f)
				d.eChan <- err

			}
		}

	}(period)

	return d
}

func (d *TimedTerminal) Display(f format.Formatter) error {
	d.fChan <- f
	return <-d.eChan
}

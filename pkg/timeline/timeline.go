package timeline

import (
	"fmt"
	"sort"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/rsrc"
)

type Event struct {
	Date    rsrc.Day
	Message string
}

type eventsT []Event

func CompileEvents(
	c charts.Charts,
	registered, from, before rsrc.Day,
) (events []Event) {

	f, ok := from.Midnight()
	if !ok {
		return
	}
	b, ok := before.Midnight()
	if !ok {
		return
	}

	entries := c.FindEntryDatesDynamic(registered, 2)

	for _, entry := range entries {
		ed, _ := entry.Date.Midnight()
		if ed >= f && ed < b {
			events = append(events, Event{
				entry.Date,
				fmt.Sprintf("enter %v", entry.Name),
			})
		}
	}

	sort.Sort(eventsT(events))

	return
}

func (es eventsT) Len() int      { return len(es) }
func (es eventsT) Swap(i, j int) { es[i], es[j] = es[j], es[i] }
func (es eventsT) Less(i, j int) bool {
	a, _ := es[i].Date.Midnight()
	b, _ := es[j].Date.Midnight()
	return a < b
}

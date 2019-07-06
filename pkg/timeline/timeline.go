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

	if from == nil || before == nil {
		return
	}

	entries := c.FindEntryDatesDynamic(2)

	for _, entry := range entries {
		ed := entry.Date.Midnight()
		if ed >= from.Midnight() && ed < before.Midnight() {
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
	return es[i].Date.Midnight() < es[j].Date.Midnight()
}

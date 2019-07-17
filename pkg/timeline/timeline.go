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

func (es eventsT) Len() int      { return len(es) }
func (es eventsT) Swap(i, j int) { es[i], es[j] = es[j], es[i] }
func (es eventsT) Less(i, j int) bool {
	return es[i].Date.Midnight() < es[j].Date.Midnight()
}

func CompileEvents(
	c charts.Charts,
	registered, from, before rsrc.Day, // TODO remove registered, rename from & before
) (events []Event) {

	if from == nil || before == nil {
		return
	}

	// TODO extract function
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

	events = append(events,
		CompileNumberOne(
			c.Fade(365),
			charts.Interval{Begin: from, Before: before})...)

	sort.Sort(eventsT(events))

	return
}

// CompileNumberOne collects the succession of holders of the number one spot
// in the charts.
func CompileNumberOne(
	cha charts.Charts,
	interval charts.Interval) (events []Event) {
	ranks := cha.Rank()

	iBegin := cha.Headers.Index(interval.Begin)
	if iBegin < 0 {
		iBegin = 0
	}

	iEnd := cha.Headers.Index(interval.Before)
	if iEnd >= cha.Len() {
		iEnd = cha.Len() - 1
	}

	if iEnd < iBegin {
		return
	}

	if iBegin == 0 {
		events = append(events, Event{
			cha.Headers.At(0).Begin,
			fmt.Sprintf("top at begin is '%v'", top(ranks, 0).FullTitle()),
		})

		iBegin++
	}

	tmpTop := ""
	tmpI := 0
	for i := 0; i < iEnd; i++ {
		topKey := top(ranks, i).FullTitle()
		if tmpTop != topKey {
			if i >= iBegin {
				tNow := cha.Headers.At(i).Begin.Time()
				tLast := cha.Headers.At(tmpI).Begin.Time()
				duration := tNow.Sub(tLast)
				diff := duration.Hours() / 24

				events = append(events, Event{
					cha.Headers.At(i).Begin,
					fmt.Sprintf("'%v' -> '%v' (%vd)",
						tmpTop, topKey, diff),
				})
			}

			tmpTop = topKey
			tmpI = i
		}
	}

	return
}

func top(
	ranks charts.Charts,
	idx int) charts.Key {

	for i, key := range ranks.Keys {
		if ranks.Values[i][idx] == 1 {
			return key
		}
	}

	return nil
}

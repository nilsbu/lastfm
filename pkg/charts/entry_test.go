package charts

import (
	"reflect"
	"testing"

	"github.com/nilsbu/lastfm/pkg/rsrc"
)

func TestChartsFindEntryDates(t *testing.T) {
	cases := []struct {
		sums       Charts
		threshold  float64
		entryDates []EntryDate
	}{
		{
			Charts{"A": []float64{2, 3, 4}},
			3,
			[]EntryDate{{"A", rsrc.ParseDay("2018-01-02")}},
		},
		{
			Charts{
				"A": []float64{2, 3, 4},
				"B": []float64{10, 10, 11}},
			10,
			[]EntryDate{{"B", rsrc.ParseDay("2018-01-01")}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			entryDates := c.sums.FindEntryDates(
				rsrc.ParseDay("2018-01-01"),
				c.threshold)

			if !reflect.DeepEqual(c.entryDates, entryDates) {
				t.Errorf("%v != %v", c.entryDates, entryDates)
			}
		})
	}
}

func TestFilterEntryDates(t *testing.T) {
	cases := []struct {
		pre    []EntryDate
		post   []EntryDate
		cutoff rsrc.Day
	}{
		{
			[]EntryDate{
				{"A", rsrc.ParseDay("2018-01-02")},
				{"B", rsrc.ParseDay("2017-01-02")}},
			[]EntryDate{{"A", rsrc.ParseDay("2018-01-02")}},
			rsrc.ParseDay("2018-01-02"),
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			post := FilterEntryDates(c.pre, c.cutoff)

			if !reflect.DeepEqual(c.post, post) {
				t.Errorf("%v != %v", c.post, post)
			}
		})
	}
}

func TestChartsGetYearPartition(t *testing.T) {
	cases := []struct {
		sums      Charts
		partition Partition
	}{
		{
			Charts{
				"A": []float64{0, 0, 0, 0, 1, 1, 1},
				"B": []float64{1, 1, 1, 1, 0, 0, 0},
				"C": []float64{0, 0, 0, 0, 0, 0, 0}},
			mapPart{
				assoc:      map[string]string{"A": "2018", "B": "2017"},
				partitions: []string{"2017", "2018", "-"},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			partition := c.sums.GetYearPartition(rsrc.ParseDay("2017-12-29"), 2)

			if !reflect.DeepEqual(c.partition, partition) { // TODO test partition API
				t.Errorf("%v != %v", c.partition, partition)
			}
		})
	}
}

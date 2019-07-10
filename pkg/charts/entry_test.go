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
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("A")},
				Values:  [][]float64{{2, 3, 4}}},
			3,
			[]EntryDate{{"A", rsrc.ParseDay("2018-01-02")}},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2018-01-01"), rsrc.ParseDay("2018-01-04")),
				Keys:    []Key{simpleKey("A"), simpleKey("B")},
				Values:  [][]float64{{2, 3, 4}, {10, 10, 11}}},
			10,
			[]EntryDate{{"B", rsrc.ParseDay("2018-01-01")}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			entryDates := c.sums.FindEntryDates(c.threshold)

			if !reflect.DeepEqual(c.entryDates, entryDates) {
				t.Errorf("%v != %v", c.entryDates, entryDates)
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
				Headers: Days(rsrc.ParseDay("2017-12-29"), rsrc.ParseDay("2017-12-29")),
				Keys:    []Key{},
				Values:  [][]float64{}},
			mapPart{
				assoc:      map[string]string{},
				partitions: []string{"2017", "-"},
			},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2017-12-29"), rsrc.ParseDay("2018-01-05")),
				Keys:    []Key{simpleKey("A"), simpleKey("B"), simpleKey("C")},
				Values: [][]float64{
					{0, 0, 0, 0, 1, 1, 1},
					{1, 1, 1, 1, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0}}},
			mapPart{
				assoc:      map[string]string{"A": "2018", "B": "2017"},
				partitions: []string{"2017", "2018", "-"},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			partition := c.sums.GetYearPartition(2)

			if !reflect.DeepEqual(c.partition, partition) { // TODO test partition API
				t.Errorf("%v != %v", c.partition, partition)
			}
		})
	}
}

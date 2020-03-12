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
				assoc:      map[string]Key{},
				partitions: []Key{tagKey("2017"), tagKey("-")},
			},
		},
		{
			Charts{
				Headers: Days(rsrc.ParseDay("2017-11-01"), rsrc.ParseDay("2018-05-01")),
				Keys:    []Key{simpleKey("A"), simpleKey("B"), simpleKey("C")},
				Values: [][]float64{
					append(repeat(0, 30+31+31), repeat(1, 28+31+30)...),
					append(repeat(1, 30+31+31), repeat(0, 28+31+30)...),
					append(repeat(0, 30+31+31), repeat(0, 28+31+30)...),
				}},
			mapPart{
				assoc:      map[string]Key{"A": tagKey("2018"), "B": tagKey("2017")},
				partitions: []Key{tagKey("2017"), tagKey("2018"), tagKey("-")},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			partition := c.sums.GetYearPartition(2)

			if !reflect.DeepEqual(c.partition.Partitions(), partition.Partitions()) {
				t.Errorf("expected partitions %v but got %v",
					c.partition.Partitions(), partition.Partitions())
			}

			for _, key := range c.sums.Keys {
				if c.partition.Get(key) != partition.Get(key) {
					t.Errorf("partition for key '%v': expected '%v' but got '%v'",
						key, c.partition.Get(key), partition.Get(key))
				}
			}
		})
	}
}

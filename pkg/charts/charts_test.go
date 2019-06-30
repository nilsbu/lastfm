package charts

import (
	"reflect"
	"sort"
	"testing"
)

func TestCompile(t *testing.T) {
	cases := []struct {
		days   []map[string][]float64
		charts Charts
	}{
		{
			[]map[string][]float64{},
			Charts{},
		},
		{
			[]map[string][]float64{{}},
			Charts{},
		},
		{
			[]map[string][]float64{
				{"ASD": []float64{2}},
				{"WASD": []float64{1}},
				{"ASD": []float64{13}, "WASD": []float64{4}},
			},
			Charts{"ASD": []float64{2, 0, 13}, "WASD": []float64{0, 1, 4}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			charts := Compile(c.days)

			if !reflect.DeepEqual(charts, c.charts) {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v", charts, c.charts)
			}
		})
	}
}

func TestChartsUnravelDays(t *testing.T) {
	cases := []struct {
		charts Charts
		days   []map[string][]float64
	}{
		{
			Charts{},
			[]map[string][]float64{},
		},
		{
			Charts{"A": []float64{}},
			[]map[string][]float64{},
		},
		{
			Charts{"ASD": []float64{2, 0, 13}, "WASD": []float64{0, 1, 4}},
			[]map[string][]float64{
				{"ASD": []float64{2}},
				{"WASD": []float64{1}},
				{"ASD": []float64{13}, "WASD": []float64{4}},
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			days := c.charts.UnravelDays()

			if !reflect.DeepEqual(days, c.days) {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v", days, c.days)
			}
		})
	}
}

func TestChartsKeys(t *testing.T) {
	cases := []struct {
		charts Charts
		keys   []string
	}{
		{
			Charts{},
			[]string{},
		},
		{
			Charts{
				"xx": []float64{32, 45},
				"yy": []float64{32, 45},
			},
			[]string{"xx", "yy"},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			keys := c.charts.Keys()

			sort.Strings(keys)
			sort.Strings(c.keys)
			if !reflect.DeepEqual(keys, c.keys) {
				t.Errorf("wrong data (sorted):\nhas:  %v\nwant: %v",
					keys, c.keys)
			}
		})
	}
}

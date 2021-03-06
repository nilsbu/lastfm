package charts

import (
	"reflect"
	"sort"
	"testing"
)

func TestColumnSort(t *testing.T) {
	cases := []struct {
		col    Column
		sorted Column
	}{
		{
			Column{},
			Column{},
		},
		{
			Column{{"a", 1}},
			Column{{"a", 1}},
		},
		{
			Column{{"a", 1}, {"b", 2}},
			Column{{"b", 2}, {"a", 1}},
		},
		{
			Column{{"c", 1}, {"b", 2}, {"a", 1}},
			Column{{"b", 2}, {"a", 1}, {"c", 1}},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			sort.Sort(c.col)
			if !reflect.DeepEqual(c.col, c.sorted) {
				t.Errorf("not sorted correctly:\nwant: %v\nhas:  %v", c.sorted, c.col)
			}
		})
	}
}

func TestColumnTop(t *testing.T) {
	testCases := []struct {
		column Column
		n      int
		top    Column
	}{
		{
			Column{},
			0,
			Column{},
		},
		{
			Column{Score{"X", 4}},
			0,
			Column{},
		},
		{
			Column{Score{"X", 3}, Score{"lol", 2}, Score{"o0o", 0}},
			2,
			Column{Score{"X", 3}, Score{"lol", 2}},
		},
		{
			Column{Score{"X", 3}, Score{"lol", 2}, Score{"o0o", 0}},
			4,
			Column{Score{"X", 3}, Score{"lol", 2}, Score{"o0o", 0}},
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			top := tc.column.Top(tc.n)

			if !reflect.DeepEqual(top, tc.top) {
				t.Errorf("wrong data:\nhas:  %v\nwant: %v", top, tc.top)
			}
		})
	}
}

func TestColumnSum(t *testing.T) {
	cases := []struct {
		col Column
		sum float64
	}{
		{Column{}, 0},
		{Column{{"a", 10}, {"b", 2.5}}, 12.5},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			sum := c.col.Sum()
			if sum != c.sum {
				t.Errorf("got %v, expected %v", sum, c.sum)
			}
		})
	}
}

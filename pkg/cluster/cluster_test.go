package cluster

import (
	"testing"

	"github.com/nilsbu/lastfm/pkg/charts"
)

type partition map[string][]string

func (p partition) Partitions() []charts.Key {
	keys := []charts.Key{}
	for k := range p {
		keys = append(keys, key(k))
	}
	return keys
}
func (p partition) Get(k charts.Key) charts.Key {
	for name, vs := range p {
		for _, v := range vs {
			if v == k.String() {
				return key(name)
			}
		}
	}
	return nil
}

func collect(similars map[string]map[string]float32) []string {
	artists := []string{}
	for key, similar := range similars {
		artists = append(artists, key)
		for artist := range similar {
			artists = append(artists, artist)
		}
	}

	return artists
}

func TestGreedy(t *testing.T) {
	for _, c := range []struct {
		descr    string
		similars map[string]map[string]float32
		weights  charts.Column
		clusters charts.Partition
	}{
		{
			"empty input",
			map[string]map[string]float32{},
			charts.Column{},
			partition{},
		},
		{
			"only one input",
			map[string]map[string]float32{
				"A": {"a": 1}},
			charts.Column{{Name: "A", Score: 100}},
			partition{
				"A": {"A"}},
		},
		{
			"A references B",
			map[string]map[string]float32{
				"A": {"a": 1, "B": .55},
				"B": {"a": 1}},
			charts.Column{
				{Name: "A", Score: 100},
				{Name: "B", Score: 56}},
			partition{
				"A": {"A", "B"}},
		},
		{
			"B references A",
			map[string]map[string]float32{
				"A": {"a": 1},
				"B": {"a": 1, "A": .7}},
			charts.Column{
				{Name: "A", Score: 100},
				{Name: "B", Score: 56}},
			partition{
				"A": {"A", "B"}},
		},
		{
			"A and B point to C",
			map[string]map[string]float32{
				"A": {"a": 1, "C": .66},
				"B": {"a": 1, "C": .7},
				"C": {"D": .4},
				"D": {"A": .2}},
			charts.Column{
				{Name: "A", Score: 100},
				{Name: "B", Score: 56},
				{Name: "C", Score: 22},
				{Name: "D", Score: 86}},
			partition{
				"A": {"A", "B", "C"},
				"D": {"D"}},
		},
		{
			"transitive",
			map[string]map[string]float32{
				"A": {"a": 1, "B": .66},
				"B": {"a": 1, "C": .7},
				"C": {"a": .4}},
			charts.Column{
				{Name: "A", Score: 100},
				{Name: "B", Score: 56},
				{Name: "C", Score: 22},
				{Name: "D", Score: 86}},
			partition{
				"A": {"A", "B", "C"}},
		},
		{
			"circle",
			map[string]map[string]float32{
				"A": {"B": .66},
				"B": {"C": .7},
				"C": {"D": .6},
				"D": {"A": .6}},
			charts.Column{
				{Name: "A", Score: 100},
				{Name: "B", Score: 56},
				{Name: "C", Score: 22},
				{Name: "D", Score: 86}},
			partition{
				"A": {"A", "B", "C", "D"}},
		},
	} {
		t.Run(c.descr, func(t *testing.T) {
			clusters := Greedy(
				c.similars,
				c.weights,
				0.5)

			if len(c.clusters.Partitions()) != len(clusters.Partitions()) {
				t.Fatalf("expected %v partitions but got %v",
					len(c.clusters.Partitions()), len(clusters.Partitions()))
			}
			for _, cluster := range c.clusters.Partitions() {
				found := false
				for _, c := range clusters.Partitions() {
					if cluster == c {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("cluster '%v' not found", cluster)
				}
			}

			for _, artist := range collect(c.similars) {
				want := c.clusters.Get(key(artist))
				has := clusters.Get(key(artist))
				if want != has {
					t.Errorf("expected '%v' in '%v' but is in '%v'",
						artist, want, has)
				}
			}
		})
	}
}

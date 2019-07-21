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
	defaultTags := []charts.Tag{{Name: "t", Weight: 100}}

	for _, c := range []struct {
		descr    string
		similars map[string]map[string]float32
		weights  charts.Column
		clusters charts.Partition
		names    map[string]string
	}{
		{
			"empty input",
			map[string]map[string]float32{},
			charts.Column{},
			partition{},
			map[string]string{},
		},
		{
			"only one input",
			map[string]map[string]float32{
				"A": {"a": 1}},
			charts.Column{{Name: "A", Score: 100}},
			partition{
				"A": {"A"}},
			map[string]string{"A": "A"},
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
			map[string]string{"A": "A, B"},
		},
		{
			"B references A",
			map[string]map[string]float32{
				"A": {"a": 1},
				"B": {"a": 1, "A": .7}},
			charts.Column{
				{Name: "A", Score: 100},
				{Name: "B", Score: 12}},
			partition{
				"A": {"A", "B"}},
			map[string]string{"A": "A+"},
		},
		{
			"A and B point to C",
			map[string]map[string]float32{
				"A": {"a": 1, "C": .66},
				"B": {"a": 1, "C": .7},
				"C": {"a": 1, "D": .005},
				"D": {"a": 1, "A": .002}},
			charts.Column{
				{Name: "A", Score: 100},
				{Name: "B", Score: 56},
				{Name: "C", Score: 22},
				{Name: "D", Score: 86}},
			partition{
				"A": {"A", "B", "C"},
				"D": {"D"}},
			map[string]string{"A": "A, B, …", "D": "D"},
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
			map[string]string{"A": "A, B, …"},
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
			map[string]string{"A": "A, D, B, …"},
		},
	} {
		t.Run(c.descr, func(t *testing.T) {
			tags := map[string][]charts.Tag{
				"A": defaultTags,
				"B": defaultTags,
				"C": defaultTags,
				"D": defaultTags,
			}
			clusters := Greedy(
				c.similars,
				tags,
				c.weights,
				0.5)

			if len(c.clusters.Partitions()) != len(clusters.Partitions()) {
				t.Fatalf("expected %v partitions but got %v (%v vs. %v)",
					len(c.clusters.Partitions()), len(clusters.Partitions()),
					c.clusters.Partitions(), clusters.Partitions())
			}
			for _, cluster := range c.clusters.Partitions() {
				found := false
				for _, c := range clusters.Partitions() {
					if cluster.String() == c.String() {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("cluster '%v' not found, has %v",
						cluster, clusters.Partitions())
				}
			}

			for _, artist := range collect(c.similars) {
				expectedKey := c.clusters.Get(key(artist))
				if expectedKey != nil {
					want := expectedKey.String()
					has := clusters.Get(key(artist)).String()
					if want != has {
						t.Errorf("expected '%v' in '%v' but is in '%v'",
							artist, want, has)
					}
				}
			}

			for mainArtist, title := range c.names {
				if gotMainArtist := clusters.Get(key(mainArtist)); gotMainArtist == nil {
					t.Fatalf("cluster '%v' missing (needed to check title)", mainArtist)
				} else {
					if title != gotMainArtist.FullTitle() {
						t.Errorf("expect title of '%v' to be '%v' but is '%v'",
							mainArtist, title, gotMainArtist.FullTitle())
					}
				}
			}
		})
	}
}

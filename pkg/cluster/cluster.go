package cluster

import "github.com/nilsbu/lastfm/pkg/charts"

type key string

func (k key) String() string {
	return string(k)
}
func (k key) ArtistName() string {
	return string(k)
}
func (k key) FullTitle() string {
	return string(k)
}

type cluster struct {
	Name    string
	Artists []string
}

type clusters []cluster

func (cs clusters) Partitions() []charts.Key {
	keys := []charts.Key{}
	for _, cluster := range cs {
		keys = append(keys, key(cluster.Name))
	}
	return keys
}
func (cs clusters) Get(k charts.Key) charts.Key {
	for _, cluster := range cs {
		for _, v := range cluster.Artists {
			if v == k.String() {
				return key(cluster.Name)
			}
		}
	}
	return nil
}

func Greedy(
	similar map[string]map[string]float32,
	weights charts.Column,
	threshold float32,
) charts.Partition {
	cs := map[string]*cluster{}
	orientation := map[string]string{}
	unused := map[string]bool{}

	for artist := range similar {
		cs[artist] = &cluster{Name: artist, Artists: []string{artist}}
		orientation[artist] = artist
		unused[artist] = true
	}

	greedy(similar, cs, orientation, unused, threshold)

	weightsMap := map[string]float64{}
	for _, weight := range weights {
		weightsMap[weight.Name] = weight.Score
	}

	outClusters := clusters{}
	for _, cluster := range cs {
		for _, a := range cluster.Artists {
			if weightsMap[cluster.Name] < weightsMap[a] {
				cluster.Name = a
			}
		}

		outClusters = append(outClusters, *cluster)
	}

	return outClusters
}

func greedy(
	similar map[string]map[string]float32,
	cs map[string]*cluster,
	orientation map[string]string,
	unused map[string]bool,
	threshold float32,
) {
	for {
		if artist, ok := getNext(unused); ok {
			unused[artist] = false
			delete(unused, artist)
			for other, match := range similar[artist] {
				if match < threshold {
					continue
				}
				if c, ok := orientation[other]; !ok {
					continue
				} else {
					prePosition := orientation[artist]

					if prePosition == c {
						continue
					}

					clusterA := cs[c]
					clusterB := cs[prePosition]
					for _, a := range clusterB.Artists {
						clusterA.Artists = append(clusterA.Artists, a)
						orientation[a] = clusterA.Name
					}

					delete(cs, prePosition)
				}
			}
		} else {
			break
		}
	}
}

func getNext(unused map[string]bool) (artist string, ok bool) {
	for name, free := range unused {
		if free {
			return name, true
		}
	}
	return "", false
}

package cluster

import (
	"sort"

	"github.com/nilsbu/lastfm/config"
	"github.com/nilsbu/lastfm/pkg/charts"
)

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
	Name    charts.Key
	Artists []string
}

type clusters []cluster

func (cs clusters) Partitions() []charts.Key {
	keys := []charts.Key{}
	for _, cluster := range cs {
		keys = append(keys, cluster.Name)
	}
	return keys
}
func (cs clusters) Get(k charts.Key) charts.Key {
	for _, cluster := range cs {
		for _, v := range cluster.Artists {
			if v == k.String() {
				return cluster.Name
			}
		}
	}
	return nil
}

func Greedy(
	similar map[string]map[string]float32,
	tags map[string][]charts.Tag,
	weights charts.Column,
	threshold float32,
) charts.Partition {

	cs := map[string]*cluster{}
	orientation := map[string]string{}
	unused := map[string]bool{}

	for artist := range similar {
		cs[artist] = &cluster{Name: key(artist), Artists: []string{artist}}
		orientation[artist] = artist
		unused[artist] = true
	}

	tagSimilar := calcTagSimilar(similar, tags)

	greedy(tagSimilar, cs, orientation, unused, threshold)

	weightsMap := map[string]float64{}
	for _, weight := range weights {
		weightsMap[weight.Name] = weight.Score
	}

	outClusters := clusters{}
	weightsByCluster := map[string]charts.Column{}
	for _, cluster := range cs {
		col := charts.Column{}
		for _, a := range cluster.Artists {
			col = append(col, charts.Score{Name: a, Score: weightsMap[a]})
			if weightsMap[cluster.Name.String()] < weightsMap[a] {
				cluster.Name = key(a)
			}
		}

		outClusters = append(outClusters, *cluster)
		sort.Sort(col)
		weightsByCluster[cluster.Name.String()] = col
	}

	outClusters = assignNames(outClusters, weightsByCluster)

	return outClusters
}

func calcTagSimilar(
	artistSimilar map[string]map[string]float32,
	artistTags map[string][]charts.Tag,
) map[string]map[string]float32 {
	return adjustSimilarByTags(
		artistTags,
		normalizeSimilar(artistSimilar, 30))
}

// func calcMeanTagSum(artistTags map[string][]charts.Tag) float32{
//
// }

func adjustSimilarByTags(
	tags map[string][]charts.Tag,
	similar map[string]map[string]float32,
) map[string]map[string]float32 {
	normalized := map[string]map[string]float32{}

	for artist, aTags := range tags {
		sum := 0
		for _, tag := range aTags {
			sum += tag.Weight
		}
		if sum == 0 {
			normalized[artist] = map[string]float32{}
		} else {
			nTags := map[string]float32{}
			for _, tag := range aTags {
				nTags[tag.Name] = float32(tag.Weight) / float32(sum)
			}

			for country := range config.Countries {
				delete(nTags, country)
			}

			normalized[artist] = nTags
		}
	}

	result := map[string]map[string]float32{}
	for artist, aSimilar := range similar {
		aTags := normalized[artist]
		autocorr := correlate(aTags, aTags)

		product := map[string]float32{}
		for bro, simScore := range aSimilar {
			if bTags, ok := normalized[bro]; ok {
				corr := correlate(aTags, bTags)
				corr /= autocorr

				product[bro] = simScore * corr
			}
		}

		result[artist] = product
	}

	return result
}

func correlate(aTags, bTags map[string]float32) (corr float32) {
	for tag, val := range aTags {
		if bVal, ok := bTags[tag]; ok {
			if bVal > val {
				corr += val * val
			} else {
				corr += val * bVal
			}
		}
	}
	return
}

func normalizeSimilar(
	similar map[string]map[string]float32,
	target float32,
) map[string]map[string]float32 {
	normalized := map[string]map[string]float32{}

	for name, sim := range similar {
		var weight float32
		for _, match := range sim {
			weight += match
		}

		if weight == 0 {
			normalized[name] = sim
		} else {
			norm := map[string]float32{}
			for n, match := range sim {
				norm[n] = match * target / weight
			}
			normalized[name] = norm
		}
	}

	return normalized
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
						orientation[a] = clusterA.Name.String()
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

package cluster

import (
	"bytes"

	"github.com/nilsbu/lastfm/pkg/charts"
)

// AssignNames chooses names for the clusters.
func assignNames(
	inClusters clusters,
	weights map[string]charts.Column,
) clusters {
	renamed := clusters{}
	for _, in := range inClusters {
		newKey := charts.NewCustomKey(
			in.Name.String(),
			in.Name.String(),
			getName(weights[in.Name.String()]))

		renamed = append(renamed, cluster{
			Name:    newKey,
			Artists: in.Artists,
		})
	}

	return renamed
}

func getName(weights charts.Column) string {
	if len(weights) == 1 {
		return weights[0].Name
	}

	total := weights.Sum()
	if weights[0].Score > 0.75*total {
		return weights[0].Name + "+"
	}

	var buffer bytes.Buffer

	count := 5
	if len(weights) < 5 {
		count = len(weights)
	}

	first := true
	var sum float64
	for i := 0; i < count; i++ {
		sum += weights[i].Score

		if first {
			first = false
		} else {
			buffer.WriteString(", ")
		}

		buffer.WriteString(weights[i].Name)

		if sum > 0.75*total {
			if sum < (1-1e-5)*total {
				buffer.WriteString(", …")
			}
			break
		}
	}

	return buffer.String()
}

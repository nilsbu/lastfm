package charts

import (
	"runtime"
)

type normalizer struct {
	chartsNode
	totals   Charts
	lineChan chan normalizerJob
}

type normalizerJob struct {
	in, out    []float64
	begin, end int
	back       chan error
}

func newNormalizer(parent Charts, totals Charts) *normalizer {
	n := &normalizer{
		chartsNode: chartsNode{parent: parent},
		totals:     Cache(totals),
		lineChan:   make(chan normalizerJob),
	}

	f := func(in, out, totals []float64, begin, end int) {
		for i, v := range in {
			if totals[i+begin] > 0 {
				out[i] = v / totals[i+begin]
			} else {
				out[i] = 0
			}
		}
	}

	workers := runtime.NumCPU()
	for i := 0; i < workers; i++ {
		go func() {
			// TODO is there a way to not query the entire length of the totals?
			totals, err := n.totals.Data([]Title{StringTitle("total")}, 0, parent.Len())
			for job := range n.lineChan {
				if err == nil {
					f(job.in, job.out, totals[0], job.begin, job.end)
				}
				job.back <- err
			}

		}()
	}

	return n
}

func Normalize(c Charts) Charts {
	return newNormalizer(c, ColumnSum(c))
}

func (c *normalizer) Data(titles []Title, begin, end int) ([][]float64, error) {
	data := make([][]float64, len(titles))
	back := make(chan error, len(titles))

	for i, title := range titles {
		out := make([]float64, end-begin)
		data[i] = out
		res, err := c.parent.Data([]Title{title}, begin, end)
		if err != nil {
			return nil, err
		}

		c.lineChan <- normalizerJob{
			in:    res[0],
			out:   out,
			begin: begin,
			end:   end,
			back:  back,
		}
	}
	for range titles {

		if err := <-back; err != nil {
			return nil, err
		}
	}

	close(back)
	return data, nil
}

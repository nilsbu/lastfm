package charts2

import (
	"runtime"
)

type normalizer struct {
	chartsNode
	totals   LazyCharts
	lineChan chan normalizerJob
}

type normalizerJob struct {
	in, out    []float64
	begin, end int
	back       chan bool
}

func newNormalizer(parent LazyCharts, totals LazyCharts) *normalizer {
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
			totals := n.totals.Data([]Title{StringTitle("total")}, 0, parent.Len())[0]
			for job := range n.lineChan {
				f(job.in, job.out, totals, job.begin, job.end)
				job.back <- true
			}

		}()
	}

	return n
}

func NormalizeColumn(c LazyCharts) LazyCharts {
	return newNormalizer(c, ColumnSum(c))
}

func NormalizeGaussian(
	c LazyCharts,
	sigma float64,
	width int,
	mirrorBegin, mirrorEnd bool) LazyCharts {

	smooth := Cache(Gaussian(c, sigma, width, mirrorBegin, mirrorEnd))
	return newNormalizer(smooth, ColumnSum(smooth))
}

func (c *normalizer) Column(titles []Title, index int) []float64 {
	data := c.Data(titles, index, index+1)
	tvm := make([]float64, len(titles))
	for i := range titles {
		tvm[i] = data[i][0]
	}
	return tvm
}

func (c *normalizer) Data(titles []Title, begin, end int) [][]float64 {
	data := make([][]float64, len(titles))
	back := make(chan bool, len(titles))

	for i, title := range titles {
		out := make([]float64, end-begin)
		data[i] = out
		c.lineChan <- normalizerJob{
			in:    c.parent.Data([]Title{title}, begin, end)[0],
			out:   out,
			begin: begin,
			end:   end,
			back:  back,
		}
	}
	for range titles {
		<-back
	}
	close(back)
	return data
}

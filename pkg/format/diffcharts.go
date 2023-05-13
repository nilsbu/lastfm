package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nilsbu/lastfm/pkg/charts"
)

type DiffCharts struct {
	Charts     []charts.DiffCharts
	Ranges     charts.Ranges
	Numbered   bool
	Precision  int
	Percentage bool
}

func convertToCharts(c *DiffCharts) *Charts {
	var charts Charts
	for _, chart := range c.Charts {
		charts.Charts = append(charts.Charts, chart)
	}
	charts.Ranges = c.Ranges
	charts.Numbered = c.Numbered
	charts.Precision = c.Precision
	charts.Percentage = c.Percentage

	return &charts
}

func (c *DiffCharts) CSV(w io.Writer, decimal string) error {
	return convertToCharts(c).CSV(w, decimal)
}

func (c *DiffCharts) Plain(w io.Writer) error {
	return convertToCharts(c).Plain(w)
}

func (c *DiffCharts) HTML(w io.Writer) error {
	return convertToCharts(c).HTML(w)
}

type diffChartData struct {
	Title     string  `json:"title"`
	Value     float64 `json:"value"`
	PrevPos   int     `json:"prevPos"`
	PrevValue float64 `json:"prevValue"`
}

type diffChart struct {
	Data []diffChartData `json:"data"`
}

type diffChartJSON struct {
	Chart     diffChart `json:"chart"`
	Precision int       `json:"precision"`
}

func convertDiffDataToJSON(c charts.DiffCharts, precision int, d *data) ([]byte, error) {
	var jsonData diffChartJSON
	jsonData.Precision = precision
	for i, title := range d.titles[0] {
		value := d.values[0][i][0]
		place, prevValue, _ := c.Previous(title)
		chartData := diffChartData{
			Title:     title.String(),
			Value:     value,
			PrevPos:   place,
			PrevValue: prevValue,
		}
		jsonData.Chart.Data = append(jsonData.Chart.Data, chartData)
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return jsonBytes, fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	return jsonBytes, nil
}

func (c *DiffCharts) JSON(w io.Writer) error {

	d, err := prep(convertToCharts(c))
	if err != nil {
		return err
	}

	bytes, err := convertDiffDataToJSON(c.Charts[0], c.Precision, d)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

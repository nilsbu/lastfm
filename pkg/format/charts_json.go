package format

import (
	"encoding/json"
	"fmt"
	"io"
)

type chartData struct {
	Title string  `json:"title"`
	Value float64 `json:"value"`
}

type chart struct {
	Data []chartData `json:"data"`
}

type chartJSON struct {
	Chart chart `json:"chart"`
}

func convertDataToJSON(d *data) ([]byte, error) {
	var jsonData chartJSON
	for i, title := range d.titles[0] {
		value := d.values[0][i][0]
		chartData := chartData{
			Title: title.String(),
			Value: value,
		}
		jsonData.Chart.Data = append(jsonData.Chart.Data, chartData)
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return jsonBytes, fmt.Errorf("failed to marshal JSON data: %v", err)
	}

	return jsonBytes, nil
}

type js map[string]interface{}

func (f *Charts) JSON(w io.Writer) error {
	d, err := prep(f)
	if err != nil {
		return err
	}

	bytes, err := convertDataToJSON(d)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

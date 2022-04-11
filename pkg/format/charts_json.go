package format

import (
	"encoding/json"
	"io"
)

type js map[string]interface{}

func (f *Charts) JSON(w io.Writer) error {
	d, err := prep(f)
	if err != nil {
		return err
	}

	obj := make(js)
	obj["type"] = "bar"
	data := []float64{}
	labels := []string{}

	titles := d.titles[0]
	values := d.values[0]
	for i := range titles {
		labels = append(labels, titles[i].String())
		data = append(data, values[i][0])
	}

	obj["data"] = js{"labels": labels, "datasets": []js{{"data": data}}}
	obj["options"] = js{"indexAxis": "y"}

	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

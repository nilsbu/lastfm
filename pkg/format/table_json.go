package format

import (
	"encoding/json"
	"io"
)

func (f *Table) JSON(w io.Writer) error {
	titles := f.Charts.Titles()
	values, err := f.Charts.Data(titles, 0, f.Charts.Len())
	if err != nil {
		return err
	}

	obj := make(js)
	obj["type"] = "line"
	data := []js{}
	labels := []string{}
	for i := 0; i < len(f.Ranges.Delims)-1; i++ {
		labels = append(labels, f.Ranges.Delims[i].String())
	}

	for i := range titles {
		elem := js{
			"label": titles[i].String(),
			"data":  values[i],
		}
		data = append(data, elem)
	}

	obj["data"] = js{"labels": labels, "datasets": data}

	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	return err
}

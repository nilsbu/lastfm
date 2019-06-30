package charts

// Charts is table of daily accumulation of plays.
type Charts map[string][]float64

// Compile builds Charts from single day plays.
func Compile(days []map[string][]float64) Charts {
	size := len(days)
	charts := make(Charts)
	for i, day := range days {
		for name, plays := range day {
			if _, ok := charts[name]; !ok {
				charts[name] = make([]float64, size)
			}
			charts[name][i] = plays[0]
		}
	}

	return charts
}

// UnravelDays takes Charts and disassembles it into single day plays. It acts
// as an inverse to Compile().
func (c Charts) UnravelDays() []map[string][]float64 {
	days := []map[string][]float64{}
	for i := 0; i < c.Len(); i++ {
		day := map[string][]float64{}

		for name, line := range c {
			if line[i] != 0 {
				day[name] = []float64{line[i]}
			}
		}

		days = append(days, day)
	}

	return days
}

func (c Charts) Len() int {
	for _, line := range c {
		return len(line)
	}

	return 0
}

// Keys returns the keys of the charts.
func (c Charts) Keys() []string {
	keys := []string{}
	for key := range c {
		keys = append(keys, key)
	}
	return keys
}

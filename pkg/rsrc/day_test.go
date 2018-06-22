package rsrc

import "testing"

func TestDate(t *testing.T) {
	cases := []struct {
		day      Day
		midnight int64
		ok       bool
	}{
		{ToDay(917740800), 917740800, true}, // same time
		{ToDay(917741200), 917740800, true}, // some hours later
		{NoDay(), 0, false},                 // no valid midnight
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			midnight, ok := c.day.Midnight()

			if !ok && c.ok {
				t.Error("midnight should be ok")
			} else if ok && !c.ok {
				t.Error("midnight should not be ok")
			}
			if ok {
				if midnight != c.midnight {
					t.Errorf("got midnight '%v', expected '%v'",
						midnight, c.midnight)
				}
			}
		})
	}
}

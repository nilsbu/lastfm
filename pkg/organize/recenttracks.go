package organize

import (
	"encoding/json"
	"errors"

	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

// LoadAllDayPlays load plays from all days since the registration of the user.
// TODO consistant naming scheme between read, write, download, load, get etc.
func LoadAllDayPlays(
	user unpack.User,
	until rsrc.Day,
	r io.AsyncReader) ([]unpack.DayPlays, error) {

	untilMdn, uOK := until.Midnight()
	registered, rOK := user.Registered.Midnight()
	if !uOK {
		return nil, errors.New("parameter 'until' is no valid Day")
	} else if !rOK {
		return nil, errors.New("user has no valid registration date")
	}
	days := int((untilMdn - registered) / 86400)
	result := make([]unpack.DayPlays, days+1)
	feedback := make(chan error)
	for i := range result {
		go func(i int) {
			date := rsrc.ToDay(int64(i*86400) + registered)
			dp, err := loadDayPlays(user.Name, date, r)
			if err == nil {
				result[i] = dp
			}
			feedback <- err
		}(i)
	}
	fail := []error{}
	for i := 0; i <= days; i++ {
		err := <-feedback
		if err != nil {
			fail = append(fail, err)
		}
	}
	if len(fail) > 0 {
		return nil, fail[0] // TODO return all errors
	}
	return result, nil
}

// LoadDayPlaysResult is the result of loadDayPlays.
type LoadDayPlaysResult struct {
	DayPlays unpack.DayPlays
	Err      error
}

func loadDayPlays(
	user rsrc.Name,
	time rsrc.Day,
	r io.AsyncReader,
) (unpack.DayPlays, error) {
	rs, err := rsrc.History(user, 1, time)
	if err != nil {
		return nil, err
	}
	dp, pages, err := loadDayPlaysPage(rs, r)
	if err != nil || pages == 1 {
		return dp, err
	}

	back := make(chan LoadDayPlaysResult)
	for p := 1; p < pages; p++ {
		go func(p rsrc.Page) {
			rs, tmpErr := rsrc.History(user, p+1, time)
			var tmpDP unpack.DayPlays
			if tmpErr == nil {
				tmpDP, _, tmpErr = loadDayPlaysPage(rs, r)
			}

			back <- LoadDayPlaysResult{tmpDP, tmpErr}
		}(rsrc.Page(p))
	}

	for p := 1; p < pages; p++ {
		dpr := <-back
		if dpr.Err != nil {
			return nil, dpr.Err
		}

		for k, v := range dpr.DayPlays {
			dp[k] += v
		}
	}
	close(back)

	return dp, err
}

func loadDayPlaysPage(rs rsrc.Resource,
	r io.AsyncReader) (dp unpack.DayPlays, pages int, err error) {
	result := <-r.Read(rs)
	if result.Err != nil {
		return nil, 0, result.Err
	}

	urt := &unpack.UserRecentTracks{}
	err = json.Unmarshal(result.Data, urt)
	if err != nil {
		return
	}

	dp = unpack.CountPlays(urt)
	pages = unpack.GetTracksPages(urt)
	return
}
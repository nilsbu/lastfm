package organize

import (
	"encoding/json"

	"github.com/nilsbu/lastfm/io"
	"github.com/nilsbu/lastfm/unpack"
)

// LoadAllDayPlays load plays from all days since the registration of the user.
// TODO consistant naming scheme between read, write, download, load, get etc.
func LoadAllDayPlays(
	user unpack.User,
	until io.Midnight,
	r io.AsyncReader) ([]unpack.DayPlays, error) {

	days := int((until - user.Registered) / 86400)
	result := make([]unpack.DayPlays, days+1)
	feedback := make(chan error)
	for i := range result {
		go func(i int) {
			date := io.Midnight(i)*86400 + user.Registered
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
	user io.Name,
	time io.Midnight,
	r io.AsyncReader) (unpack.DayPlays, error) {
	dp, pages, err := loadDayPlaysPage(io.NewUserRecentTracks(user, 1, time), r)
	if err != nil || pages == 1 {
		return dp, err
	}

	back := make(chan LoadDayPlaysResult)
	for p := 1; p < pages; p++ {
		go func(p io.Page) {

			tmpDP, _, tmpErr := loadDayPlaysPage(
				io.NewUserRecentTracks(user, p+1, time), r)

			back <- LoadDayPlaysResult{tmpDP, tmpErr}
		}(io.Page(p))
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

func loadDayPlaysPage(rsrc *io.Resource,
	r io.AsyncReader) (dp unpack.DayPlays, pages int, err error) {
	result := <-r.Read(rsrc)
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

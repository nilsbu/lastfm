package organize

import (
	"github.com/nilsbu/lastfm/io"
	"github.com/nilsbu/lastfm/unpack"
)

// LoadDayPlaysResult is the result of LoadDayPlays.
type LoadDayPlaysResult struct {
	DayPlays unpack.DayPlays
	Err      error
}

// LoadDayPlays load plays from one day.
// TODO consistant naming scheme between read, write, download, load, get etc.
func LoadDayPlays(
	user io.Name,
	time io.Midnight,
	r io.AsyncReader) <-chan LoadDayPlaysResult {

	out := make(chan LoadDayPlaysResult)
	go func() {

		dp, pages, err := loadDayPlaysPage(io.NewUserRecentTracks(user, 1, time), r)
		if err != nil || pages == 1 {
			out <- LoadDayPlaysResult{dp, err}
			close(out)
			return
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
				out <- LoadDayPlaysResult{nil, dpr.Err}
				close(out)
				return
			}

			for k, v := range dpr.DayPlays {
				dp[k] += v
			}
		}
		close(back)

		out <- LoadDayPlaysResult{dp, err}
		close(out)
	}()
	return out
}

func loadDayPlaysPage(rsrc *io.Resource,
	r io.AsyncReader) (dp unpack.DayPlays, pages int, err error) {
	result := <-r.Read(rsrc)
	if result.Err != nil {
		return nil, 0, result.Err
	}

	urt, err := unpack.UnmarshalUserRecentTracks(result.Data)
	if err != nil {
		return
	}

	dp = unpack.CountPlays(urt)
	pages = unpack.GetTracksPages(urt)
	return
}

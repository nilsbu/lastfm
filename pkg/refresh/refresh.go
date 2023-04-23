package refresh

import (
	"sync"
	"time"

	"github.com/nilsbu/lastfm/pkg/charts"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type Trigger interface {
	Refresh()
}

type refreshPipeline struct {
	mtx      sync.RWMutex
	store    io.Store
	session  *unpack.SessionInfo
	pipeline pipeline.Pipeline
}

type refreshTrigger struct {
	pipeline *refreshPipeline
}

func (t *refreshTrigger) Refresh() {
	t.pipeline.mtx.Lock()
	defer t.pipeline.mtx.Unlock()
	t.pipeline.pipeline = pipeline.New(t.pipeline.session, t.pipeline.store)
}

func (rp *refreshPipeline) Execute(steps []string) (charts.Charts, error) {
	rp.mtx.RLock()
	defer rp.mtx.RUnlock()
	return rp.pipeline.Execute(steps)
}

func (rp *refreshPipeline) Registered() rsrc.Day {
	rp.mtx.RLock()
	defer rp.mtx.RUnlock()
	return rp.pipeline.Registered()
}

func (rp *refreshPipeline) Session() *unpack.SessionInfo {
	rp.mtx.RLock()
	defer rp.mtx.RUnlock()
	return rp.pipeline.Session()
}

func WrapRefresh(s io.Store, pl pipeline.Pipeline, session *unpack.SessionInfo) (pipeline.Pipeline, Trigger) {
	rpl := &refreshPipeline{
		store:    s,
		session:  session,
		pipeline: pl,
	}
	return rpl, &refreshTrigger{pipeline: rpl}
}

func PeriodicRefresh(trigger Trigger, h, m, s int) {
	for {
		// Get the current time
		now := time.Now()

		// Calculate the duration until the specified time
		targetTime := time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, now.Location())
		duration := targetTime.Sub(now)

		// If the duration is negative, add 24 hours to get the time for the next day
		if duration < 0 {
			duration += 24 * time.Hour
		}

		// Sleep until the specified time
		time.Sleep(duration)

		// Call the Refresh method on the provided Trigger object
		trigger.Refresh()
	}
}

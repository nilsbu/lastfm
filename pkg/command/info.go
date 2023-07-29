package command

import (
	"strconv"
	"strings"

	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/io"
	"github.com/nilsbu/lastfm/pkg/pipeline"
	"github.com/nilsbu/lastfm/pkg/rsrc"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type infoT struct {
	rsrc  string
	param string
}

func (cmd infoT) Execute(
	session *unpack.SessionInfo, s io.Store, pl pipeline.Pipeline, d display.Display) error {

	var f rsrc.Locator = nil
	switch cmd.rsrc {
	case "user.getInfo":
		f = rsrc.UserInfo(cmd.param)
	case "user.getRecentTracks":
		strs := strings.Split(cmd.param, " ")
		page, err := strconv.Atoi(strs[1])
		if err != nil {
			return err
		}
		day := rsrc.ParseDay(strs[2])
		f = rsrc.History(strs[0], page, day)
	case "artist.getInfo":
		f = rsrc.ArtistInfo(cmd.param)
	case "artist.getTopTags":
		f = rsrc.ArtistTags(cmd.param)
	case "tag.getInfo":
		f = rsrc.TagInfo(cmd.param)
	case "track.getInfo":
		strs := strings.Split(cmd.param, "^")
		f = rsrc.TrackInfo(strs[0], strs[1])
	}

	if msg, err := f.Path(); err != nil {
		return err
	} else {
		d.Display(&format.Message{Msg: msg})
	}

	key, err := unpack.LoadAPIKey(io.FileIO{})
	if err != nil {
		return err
	}
	if msg, err := f.URL(key); err != nil {
		return err
	} else {
		d.Display(&format.Message{Msg: msg})
	}

	return nil
}

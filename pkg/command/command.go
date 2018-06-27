package command

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/organize"
	"github.com/nilsbu/lastfm/pkg/store"
)

// Execute executes the command described in the arguments.
func Execute(args []string, sid organize.SessionID, s store.Store) error {
	cmd, err := resolve(args, sid)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = cmd.Execute(s)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

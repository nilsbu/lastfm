package command

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/io"
)

// Execute executes the command described in the arguments.
func Execute(args []string, store io.Store) error {
	cmd, err := resolve(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = cmd.Execute(store)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

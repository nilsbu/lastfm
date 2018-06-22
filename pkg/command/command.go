package command

import (
	"fmt"

	"github.com/nilsbu/lastfm/pkg/io"
)

// Execute executes the command described in the arguments.
func Execute(args []string, ioPool io.Pool) error {
	cmd, err := resolve(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = cmd.Execute(ioPool)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

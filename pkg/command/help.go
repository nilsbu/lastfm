package command

import (
	"github.com/nilsbu/lastfm/pkg/display"
	"github.com/nilsbu/lastfm/pkg/format"
	"github.com/nilsbu/lastfm/pkg/store"
	"github.com/nilsbu/lastfm/pkg/unpack"
)

type help struct{}

func (help) Execute(session *unpack.SessionInfo, s store.Store, d display.Display) error {
	for _, str := range listCommands() {
		d.Display(&format.Message{Msg: str})
	}

	return nil
}

func listCommands() []string {
	childStrs := []string{}

	for childName, childNode := range cmdRoot.nodes {
		if childName != "lastfm" {
			continue
		}

		for _, cstr := range listCommandsTree(childNode) {
			str := childName + " " + cstr
			childStrs = append(childStrs, str)
		}
	}

	return childStrs
}

func listCommandsTree(node node) []string {

	childStrs := []string{}

	if node.cmd != nil {
		childStrs = append(childStrs, getCommandDescription(node.cmd))
	}

	for childName, childNode := range node.nodes {
		for _, cstr := range listCommandsTree(childNode) {
			str := childName + " " + cstr
			childStrs = append(childStrs, str)
		}
	}

	return childStrs
}

func getCommandDescription(cmd *cmd) string {
	str := ""
	for _, param := range cmd.params {
		str += "[" + param.name + ":" + param.kind + "] "
	}

	return str + "- " + cmd.descr
}

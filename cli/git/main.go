package git

import (
	"gitlab.768bit.com/pub/vpkg"
	"gitlab.768bit.com/pub/vpkg/cli/support"
	"gopkg.in/urfave/cli.v1"
)

var VDATA *vpkg.VersionData

var GitCommand = cli.Command{
	Name:  "git",
	Usage: "Simple git operations for version.json items",
	Before: func(context *cli.Context) error {

		VDATA = support.GetVersionData()
		return nil

	},
	Subcommands: []cli.Command{
		GitTagCommand,
	},
}

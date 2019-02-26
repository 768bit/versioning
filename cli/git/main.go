package git

import (
	"github.com/768bit/vpkg"
	"github.com/768bit/vpkg/cli/support"
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

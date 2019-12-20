package main

import (
	"github.com/768bit/vpkg/cli/support"
	"gopkg.in/urfave/cli.v1"
	"log"
)

var SetVersionCommmand = cli.Command{
	Name:  "set",
	Usage: "Set Version Data Information from version.json",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "major",
			Usage: "Set major version",
		},
		cli.IntFlag{
			Name:  "minor",
			Usage: "Set minor version",
		},
		cli.IntFlag{
			Name:  "revision",
			Usage: "Set revision",
		},
	},
	Action: func(c *cli.Context) error {
		VDATA := support.GetVersionData()

		prevVersion := VDATA.FullVersionString()

		major := c.Int("major")
		minor := c.Int("minor")
		revision := c.Int("revision")

		changed := false

		if major >= 0 && major > VDATA.Major {
			VDATA.Major = major
			changed = true
		}

		if minor >= 0 && minor > VDATA.Minor {
			VDATA.Minor = minor
			changed = true
		}

		if revision >= 0 && revision > VDATA.Revision {
			VDATA.Revision = revision
			changed = true
		}

		if changed {
			log.Println("Setting Version from", prevVersion, " -> ", VDATA.FullVersionString())
			return VDATA.Save(ROOT)
		}

		return nil

	},
}

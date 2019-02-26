package main

import (
	"errors"
	"fmt"
	"gitlab.768bit.com/pub/vpkg"
	"gitlab.768bit.com/pub/vpkg/cli/support"
	"gopkg.in/urfave/cli.v1"
)

//make a new version - this will increment/set version numbers based on provided flags...

var NewVersionCommand = cli.Command{
	Name:  "new",
	Usage: "Make a new version for the version.json file: major, minor, revision, build and os package revision",
	Subcommands: []cli.Command{
		NewMajorVersionCommand,
		NewMinorVersionCommand,
		NewRevisionVersionCommand,
		NewBuildVersionCommand,
		NewOSPackageVersionCommand,
	},
}

var NewMajorVersionCommand = cli.Command{
	Name:  "major",
	Usage: "Make a new Major version",
	Action: func(c *cli.Context) error {
		VDATA := support.GetVersionData()
		VDATA.NewMajorVersion()
		return VDATA.Save(ROOT)
	},
}

var NewMinorVersionCommand = cli.Command{
	Name:  "minor",
	Usage: "Make a new Minor version",
	Action: func(c *cli.Context) error {
		VDATA := support.GetVersionData()
		VDATA.NewMinorVersion()
		return VDATA.Save(ROOT)
	},
}

var NewRevisionVersionCommand = cli.Command{
	Name:  "revision",
	Usage: "Make a new Revision",
	Action: func(c *cli.Context) error {
		VDATA := support.GetVersionData()
		VDATA.NewRevision()
		return VDATA.Save(ROOT)
	},
}

var NewBuildVersionCommand = cli.Command{
	Name:  "build",
	Usage: "Make a new Build",
	Action: func(c *cli.Context) error {
		VDATA := support.GetVersionData()
		VDATA.NewBuild()
		return VDATA.Save(ROOT)
	},
}

var NewOSPackageVersionCommand = cli.Command{
	Name:  "pkg",
	Usage: "Make a new revision of the package for a specified operating system",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "os, o",
			Usage: "OS to make a package revision for (debian, ubuntu, macos, windows)",
		},
	},
	Action: func(c *cli.Context) error {
		VDATA := support.GetVersionData()

		os := c.String("os")

		allowed, os := vpkg.CheckAllowedOs(os)

		if !allowed || os == "" {

			return errors.New(fmt.Sprintf("Requested OS %s is invalid", os))

		} else {

			VDATA.NewPkgRevision(os)

		}
		return VDATA.Save(ROOT)
	},
}

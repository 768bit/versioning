package main

import (
	"errors"
	"fmt"
	"github.com/768bit/verman"
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
		VDATA.NewMajorVersion()
		return VDATA.Save()
	},
}

var NewMinorVersionCommand = cli.Command{
	Name:  "minor",
	Usage: "Make a new Minor version",
	Action: func(c *cli.Context) error {
		VDATA.NewMinorVersion()
		return VDATA.Save()
	},
}

var NewRevisionVersionCommand = cli.Command{
	Name:  "revision",
	Usage: "Make a new Revision",
	Action: func(c *cli.Context) error {
		VDATA.NewRevision()
		return VDATA.Save()
	},
}

var NewBuildVersionCommand = cli.Command{
	Name:  "build",
	Usage: "Make a new Build",
	Action: func(c *cli.Context) error {
		VDATA.NewBuild()
		return VDATA.Save()
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

		os := c.String("os")

		allowed, os := verman.CheckAllowedOs(os)

		if !allowed || os == "" {

			return errors.New(fmt.Sprintf("Requested OS %s is invalid", os))

		} else {

			VDATA.NewPkgRevision(os)

		}
		return VDATA.Save()
	},
}

package main

import (
	"fmt"
	"github.com/768bit/verman"
  "github.com/768bit/verman/cli/git"
  "gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"strings"
)

var (
	Version   string
	Build     string
	BuildDate string
	GitCommit string
)

var VDATA *verman.VersionData

func GetVersionData() *verman.VersionData {

	return VDATA

}

func InitialiseVersionData() error {

	vd, err := verman.LoadVersionData()
	if err != nil {
		fmt.Println(err)
		vd = verman.NewVersionData()
		err = vd.Save()
		if err != nil {
			return err
		}
		fmt.Println("Created a new version.json file.")
		return InitialiseVersionData()
	}
	VDATA = vd
	return nil

}

func main() {

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "Show the version of verman",
	}

	cli.HelpFlag = cli.BoolFlag{Name: "help"}

	app := cli.NewApp()

	app.Name = "verman"

	app.Description = "Verman manages your projects version.json file. Git is supported. Automatically tag your release. Automated release versioning."

	app.Version = fmt.Sprintf("%s  Git Commit: %s  Build Date: %s", Version, GitCommit, strings.Replace(BuildDate, "_", " ", -1))

	app.Flags = []cli.Flag{}

	app.Before = func(context *cli.Context) error {

		return InitialiseVersionData()

	}

	app.Commands = []cli.Command{
		NewVersionCommand,
		GetVersionCommand,
		git.GitCommand,

	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"fmt"
	"github.com/768bit/vpkg/cli/git"
	"github.com/768bit/vpkg/cli/support"
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

var ROOT string = ""

func main() {

	cwd, _ := os.Getwd()

	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "Show the version of verman",
	}

	cli.HelpFlag = cli.BoolFlag{Name: "help"}

	app := cli.NewApp()

	app.Name = "verman"

	app.Description = "Verman manages your projects version.json file. Git is supported. Automatically tag your release. Automated release versioning."

	app.Version = fmt.Sprintf("%s  Git Commit: %s  Build Date: %s", Version, GitCommit, strings.Replace(BuildDate, "_", " ", -1))

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "target-folder, t",
			Usage: "The target folder for git operations etc when using a namespaced version.json file",
			Value: ".",
		},
		cli.StringFlag{
			Name:  "namespace, n",
			Usage: "The namespace to use",
			Value: "",
		},
		cli.StringFlag{
			Name:  "root",
			Usage: "The root folder to use",
			Value: cwd,
		},
	}

	app.Before = func(context *cli.Context) error {

		ROOT = context.String("root")

		return support.InitialiseVersionData(context.String("namespace"), context.String("target-folder"), context.String("root"))

	}

	app.Commands = []cli.Command{
		NewVersionCommand,
		GetVersionCommand,
		SetVersionCommmand,
		git.GitCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

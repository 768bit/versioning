package git

import (
  "fmt"
  "github.com/768bit/verman"
  "gopkg.in/urfave/cli.v1"
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

var GitCommand = cli.Command{
  Name:    "git",
  Usage:   "Simple git operations for version.json items",
  Before : func(context *cli.Context) error {

    return InitialiseVersionData()

  },
  Subcommands: []cli.Command{
    GitTagCommand,
  },
}

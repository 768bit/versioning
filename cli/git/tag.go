package git

import "gopkg.in/urfave/cli.v1"

var GitTagCommand = cli.Command{
  Name:  "tag",
  Usage: "Tag the git repository with the appropriate version data.",
  Flags: []cli.Flag{
    cli.StringFlag{
      Name:  "message, m",
      Usage: "Tage message to use",
    },
    cli.BoolFlag{
      Name: "commit, c",
      Usage: "Use the commit flag to commit the pero at the same time.",
    },
  },
  Action: func(c *cli.Context) error {

    if c.Bool("commit") {

      //commit now...

      if err := VDATA.PerformGitCommit(c.String("message")); err != nil {
        return err
      }

    }

    return VDATA.TagGitCommit(c.String("message"))
  },
}

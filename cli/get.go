package main

import (
	"errors"
	"fmt"
	"github.com/768bit/verman"
	"gopkg.in/urfave/cli.v1"
)

var GetVersionCommand = cli.Command{
	Name:  "get",
	Usage: "Get Version Data Information from version.json",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "json",
			Usage: "Return the version data as JSON",
		},
		cli.BoolFlag{
			Name:  "raw, r",
			Usage: "Return just the string with no other text",
		},
		cli.BoolFlag{
			Name:  "full, f",
			Usage: "Return the full version string",
		},
		cli.BoolFlag{
			Name:  "build, b",
			Usage: "Return the full version string including build id",
		},
		cli.BoolFlag{
			Name:  "date, d",
			Usage: "Return the full version string including build id",
		},
		cli.BoolFlag{
			Name:  "uuid, u",
			Usage: "Return the build uuid",
		},
		cli.StringFlag{
			Name:  "os, o",
			Usage: "Return the full version string for an OS package (implies -f and -b)",
		},
	},
	Action: func(c *cli.Context) error {

		if c.Bool("json") {

			//output the json now...

			if vdJson, err := VDATA.ToJSON(); err != nil {

				return err

			} else {

				fmt.Println(vdJson)

			}

			if !c.Bool("raw") {

				fmt.Print("\n")

			}

		} else if os := c.String("os"); os != "" {

			allowed, os := verman.CheckAllowedOs(os)

			if !allowed {

				return errors.New(fmt.Sprintf("Requested OS %s is invalid", os))

			} else if _, err := VDATA.GetPkgRevision(os); err != nil {

				return err

			} else if pkgVerAllowed, pkgVer := verman.PkgOsVersionString(os, VDATA); !pkgVerAllowed {

				return errors.New(fmt.Sprintf("Unable to obtain package version for os %s", os))

			} else {

				if !c.Bool("raw") {

					fmt.Print("Version: ")

				}

				fmt.Println(pkgVer)

				if !c.Bool("raw") {

					fmt.Print("\n")

				}

			}

		} else if c.Bool("build") {

			if c.Bool("full") {

				if !c.Bool("raw") {

					fmt.Print("Version: ")

				}

				fmt.Printf(VDATA.FullVersionString())

				if !c.Bool("raw") {

					fmt.Print("\n")

				}

			} else {

				if !c.Bool("raw") {

					fmt.Print("Build: ")

				}

				fmt.Println(VDATA.ShortID)

				if !c.Bool("raw") {

					fmt.Print("\n")

				}

			}

		} else if c.Bool("full") {

			if !c.Bool("raw") {

				fmt.Print("Version: ")

			}

			fmt.Printf("%d.%d.%d", VDATA.Major, VDATA.Minor, VDATA.Revision)

			if !c.Bool("raw") {

				fmt.Print("\n")

			}

		} else if c.Bool("uuid") {

			if !c.Bool("raw") {

				fmt.Print("UUID: ")

			}

			fmt.Printf("%s", VDATA.UUID)

			if !c.Bool("raw") {

				fmt.Print("\n")

			}

		} else if c.Bool("date") {

			if !c.Bool("raw") {

				fmt.Print("Date: ")

			}

			fmt.Printf("%s", VDATA.DateString)

			if !c.Bool("raw") {

				fmt.Print("\n")

			}

		}

		return nil

	},
}

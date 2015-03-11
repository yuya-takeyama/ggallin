package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandInstall,
	commandInit,
	commandBuild,
	commandRelease,
}

var commandInstall = cli.Command{
	Name:  "install",
	Usage: "",
	Description: `
`,
	Action: doInstall,
}

var commandInit = cli.Command{
	Name:  "init",
	Usage: "",
	Description: `
`,
	Action: doInit,
}

var commandBuild = cli.Command{
	Name:  "build",
	Usage: "",
	Description: `
`,
	Action: doBuild,
}

var commandRelease = cli.Command{
	Name:  "release",
	Usage: "",
	Description: `
`,
	Action: doRelease,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func doInstall(c *cli.Context) {
}

func doInit(c *cli.Context) {
}

func doBuild(c *cli.Context) {
}

func doRelease(c *cli.Context) {
}

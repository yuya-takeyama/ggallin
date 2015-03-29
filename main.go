package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "ggallin"
	app.Version = Version
	app.Usage = "Compose gox and ghr as all-in-one"
	app.Author = "Yuya Takeyama"
	app.Email = "sign.of.the.wolf.pentagram@gmail.com"
	app.Commands = Commands

	cli.VersionPrinter = printVersion

	app.Run(os.Args)
}

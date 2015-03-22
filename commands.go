package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codegangsta/cli"
)

var Commands = []cli.Command{
	commandInit,
	commandBuild,
	commandRelease,
}

var commandInit = cli.Command{
	Name:  "init",
	Usage: "Create new project",
	Description: `
`,
	Action: doInit,
}

var commandBuild = cli.Command{
	Name:  "build",
	Usage: "Build package",
	Description: `
`,
	Action: doBuild,
}

var commandRelease = cli.Command{
	Name:  "release",
	Usage: "build package and release it",
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

func doInit(c *cli.Context) {
}

func doBuild(c *cli.Context) {
	version, err := readVersion()
	panicIf(err)

	compile(version)
}

func doRelease(c *cli.Context) {
}

func compile(version string) {
	buildDir := filepath.Join("build", version)

	err := os.RemoveAll(buildDir)
	panicIf(err)

	cmd := exec.Command("gox", "-output", filepath.Join(buildDir, "{{.OS}}_{{.Arch}}", "{{.Dir}}"))
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	panicIf(err)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

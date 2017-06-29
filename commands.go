package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
)

type Build struct {
	Target  string
	Version string
	Dir     string
}

var archFlag = cli.StringFlag{
	Name:  "arch",
	Usage: "Space-separated list of architectures to build for",
}

var osFlag = cli.StringFlag{
	Name:  "os",
	Usage: "Space-separated list of operating systems to build for",
}

var osarchFlag = cli.StringFlag{
	Name:  "osarch",
	Usage: "Space-separated list of os/arch pairs to build for",
}

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
	Flags: []cli.Flag{
		archFlag,
		osFlag,
		osarchFlag,
	},
}

var commandRelease = cli.Command{
	Name:  "release",
	Usage: "build package and release it",
	Description: `
`,
	Action: doRelease,
	Flags: []cli.Flag{
		archFlag,
		osFlag,
		osarchFlag,
		cli.BoolFlag{
			Name:  "replace",
			Usage: "Replace asset if target is already exists",
		},
		cli.StringFlag{
			Name:  "username, u",
			Usage: "GitHub username",
		},
	},
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

	buildDir := filepath.Join("build", version)
	pkgDir := filepath.Join("pkg", version)

	compile(c, version, buildDir)
	pkg(version, buildDir, pkgDir)
}

func doRelease(c *cli.Context) {
	version, err := readVersion()
	panicIf(err)

	buildDir := filepath.Join("build", version)
	pkgDir := filepath.Join("pkg", version)

	compile(c, version, buildDir)
	pkg(version, buildDir, pkgDir)
	gitPush()
	release(c, version, pkgDir)
}

func compile(c *cli.Context, version, buildDir string) {
	err := os.RemoveAll(buildDir)
	panicIf(err)

	gitCommit, err := getCommitHash()
	panicIf(err)

	options := []string{
		"-ldflags",
		"-X main.GitCommit=" + gitCommit,
		"-output",
		filepath.Join(buildDir, "{{.OS}}_{{.Arch}}", "{{.Dir}}"),
	}

	if c.String("arch") != "" {
		options = append(options, "-arch", c.String("arch"))
	}
	if c.String("os") != "" {
		options = append(options, "-os", c.String("os"))
	}
	if c.String("osarch") != "" {
		options = append(options, "-osarch", c.String("osarch"))
	}

	cmd := exec.Command("gox", options...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	panicIf(err)
}

func getCommitHash() (string, error) {
	out := new(bytes.Buffer)

	cmd := exec.Command("git", "describe", "--always")
	cmd.Stdout = out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func pkg(version, buildDir, pkgDir string) {
	err := os.RemoveAll(pkgDir)
	panicIf(err)

	err = os.MkdirAll(pkgDir, 0755)
	panicIf(err)

	dirs, err := ioutil.ReadDir(buildDir)
	panicIf(err)

	wg := new(sync.WaitGroup)

	for _, dir := range dirs {
		wg.Add(1)

		go func(dir os.FileInfo, pkgDir string, wg *sync.WaitGroup) {
			build := &Build{
				Target:  dir.Name(),
				Version: version,
				Dir:     filepath.Join(buildDir, dir.Name()),
			}

			makeZip(build, pkgDir, wg)
		}(dir, pkgDir, wg)
	}

	wg.Wait()

	fmt.Fprintf(os.Stderr, "Package files created into %s\n", pkgDir)
}

func makeZip(build *Build, pkgDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	fileInfos, err := ioutil.ReadDir(build.Dir)
	panicIf(err)

	zipFile, err := os.Create(filepath.Join(pkgDir, build.Target+"_"+build.Version+".zip"))
	panicIf(err)

	zipWriter := zip.NewWriter(zipFile)

	for _, fileInfo := range fileInfos {
		file, err := os.Open(filepath.Join(build.Dir, fileInfo.Name()))
		panicIf(err)

		fileHeader, err := zip.FileInfoHeader(fileInfo)
		panicIf(err)

		f, err := zipWriter.CreateHeader(fileHeader)
		panicIf(err)

		_, err = io.Copy(f, file)
		panicIf(err)
	}

	err = zipWriter.Close()
	panicIf(err)
}

func gitPush() {
	cmd := exec.Command("git", "push", "origin", "master")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	panicIf(err)
}

func release(c *cli.Context, version, pkgDir string) {
	options := []string{}

	username := c.String("username")
	if username != "" {
		options = append(options, "--username", username)
	}

	if c.Bool("replace") {
		options = append(options, "--replace")
	}

	options = append(options, "v"+version, pkgDir)

	cmd := exec.Command("ghr", options...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	panicIf(err)
}

func printVersion(c *cli.Context) {
	fmt.Printf("%s v%s, build %s\n", c.App.Name, Version, GitCommit)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

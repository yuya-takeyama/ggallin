package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/codegangsta/cli"
)

type Build struct {
	Target  string
	Version string
	Dir     string
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

	buildDir := filepath.Join("build", version)
	pkgDir := filepath.Join("pkg", version)

	compile(version, buildDir)
	pkg(version, buildDir, pkgDir)
}

func doRelease(c *cli.Context) {
	version, err := readVersion()
	panicIf(err)

	buildDir := filepath.Join("build", version)
	pkgDir := filepath.Join("pkg", version)

	compile(version, buildDir)
	pkg(version, buildDir, pkgDir)
	release(version, pkgDir)
}

func compile(version, buildDir string) {
	err := os.RemoveAll(buildDir)
	panicIf(err)

	cmd := exec.Command("gox", "-output", filepath.Join(buildDir, "{{.OS}}_{{.Arch}}", "{{.Dir}}"))
	err = cmd.Run()
	panicIf(err)
}

func pkg(version, buildDir, pkgDir string) {
	err := os.RemoveAll(pkgDir)
	panicIf(err)

	err = os.MkdirAll(pkgDir, 0755)
	panicIf(err)

	buildCh := make(chan *Build)
	quitCh := make(chan bool)

	dirs, err := ioutil.ReadDir(buildDir)
	panicIf(err)

	wg := new(sync.WaitGroup)

	go zipper(buildCh, quitCh, pkgDir, wg)

	for _, dir := range dirs {
		wg.Add(1)

		build := &Build{
			Target:  dir.Name(),
			Version: version,
			Dir:     filepath.Join(buildDir, dir.Name()),
		}
		buildCh <- build
	}

	wg.Wait()
	quitCh <- true

	fmt.Fprintf(os.Stderr, "Package files created into %s\n", pkgDir)
}

func zipper(buildCh chan *Build, quitCh chan bool, pkgDir string, wg *sync.WaitGroup) {
	for {
		select {
		case build := <-buildCh:
			makeZip(build, pkgDir, wg)

		case <-quitCh:
			return
		}
	}
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

func release(version, pkgDir string) {
	cmd := exec.Command("ghr", "v"+version, pkgDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	panicIf(err)
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

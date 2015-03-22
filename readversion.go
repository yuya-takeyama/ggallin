package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func readVersion() (string, error) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}

	srcVersionFile, err := os.Open("version.go")
	if err != nil {
		return "", err
	}

	destVersionFilePath := path.Join(dir, "version.go")
	destVersionFile, err := os.Create(destVersionFilePath)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(destVersionFile, srcVersionFile)
	if err != nil {
		return "", err
	}

	mainFilePath := path.Join(dir, "main.go")
	mainFile, err := os.Create(mainFilePath)
	if err != nil {
		return "", err
	}

	mainCode := []byte(`package main

import (
	"fmt"
)

func main() {
	fmt.Print(Version)
}

`)
	_, err = mainFile.Write(mainCode)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	cmd := exec.Command("go", "run", destVersionFilePath, mainFilePath)
	cmd.Stdout = buf

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

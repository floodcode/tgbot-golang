package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxRunTime = 8 * time.Second
	progName   = "main.go"
)

func runCode(src string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "sandbox")
	if err != nil {
		return "", fmt.Errorf("error creating temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Generate imports
	in := filepath.Join(tmpDir, progName)
	cmd := exec.Command("goimports")
	inPipe, _ := cmd.StdinPipe()
	inPipe.Write([]byte(src))
	inPipe.Close()
	outImports := bytes.Buffer{}
	cmd.Stdout = &outImports
	cmd.Start()
	cmd.Wait()

	if err := ioutil.WriteFile(in, outImports.Bytes(), 0400); err != nil {
		return "", fmt.Errorf("error creating temp file %q: %v", in, err)
	}

	// Check package name
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, in, nil, parser.PackageClauseOnly)
	if err == nil && f.Name.Name != "main" {
		return "", errors.New("package name must be main")
	}

	// Build source
	exe := filepath.Join(tmpDir, "a.out")
	cmd = exec.Command("go", "build", "-o", exe, in)
	cmd.Env = []string{"GOOS=nacl", "GOARCH=amd64p32", "GOPATH=" + os.Getenv("GOPATH")}
	if out, err := cmd.CombinedOutput(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			errs := strings.Replace(string(out), in, progName, -1)
			errs = strings.Replace(errs, "# command-line-arguments\n", "", 1)
			return "", errors.New(errs)
		}

		return "", fmt.Errorf("error building go source: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxRunTime)
	defer cancel()

	// Execute built program under sandbox
	cmd = exec.CommandContext(ctx, "./sel_ldr_x86_64", "-l", "/dev/null", "-S", "-e", exe)
	rec := new(Recorder)
	cmd.Stdout = rec.Stdout()
	cmd.Stderr = rec.Stderr()
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", errors.New("process took too long")
		}
		if _, ok := err.(*exec.ExitError); !ok {
			return "", fmt.Errorf("error running sandbox: %v", err)
		}
	}

	// Decode program output
	events, err := rec.Events()
	if err != nil {
		return "", fmt.Errorf("error decoding events: %v", err)
	}

	var outputString string
	for _, event := range events {
		outputString += event.Message
	}

	return outputString, nil
}

func formatCode(src string) string {
	cmd := exec.Command("gofmt")
	inPipe, _ := cmd.StdinPipe()
	inPipe.Write([]byte(src))
	inPipe.Close()
	outImports := bytes.Buffer{}
	cmd.Stdout = &outImports
	cmd.Start()
	cmd.Wait()

	return string(outImports.Bytes())
}

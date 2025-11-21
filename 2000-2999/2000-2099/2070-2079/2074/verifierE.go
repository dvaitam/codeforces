package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSolutionPath = "2000-2999/2000-2099/2070-2079/2074/2074E.go"

func runProgram(path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader("")
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout running %s", path)
		}
		return "", fmt.Errorf("failed to run %s: %v\n%s", path, err, errBuf.String())
	}
	return strings.TrimSpace(outBuf.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/2074E_binary")
		os.Exit(1)
	}
	candidatePath := os.Args[1]

	refPath := refSolutionPath
	if !filepath.IsAbs(refPath) {
		refPath, _ = filepath.Abs(refPath)
	}
	candAbs := candidatePath
	if !filepath.IsAbs(candAbs) {
		candAbs, _ = filepath.Abs(candAbs)
	}

	expected, err := runProgram(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to run reference solution: %v\n", err)
		os.Exit(1)
	}

	actual, err := runProgram(candAbs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}

	if strings.TrimSpace(actual) != expected {
		fmt.Fprintf(os.Stderr, "output mismatch.\nexpected: %q\ngot:      %q\n", expected, strings.TrimSpace(actual))
		os.Exit(1)
	}

	fmt.Println("All tests passed")
}

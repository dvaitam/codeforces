package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSourceD1 = "1000-1999/1300-1399/1350-1359/1357/1357D1.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/candidate")
		os.Exit(1)
	}

	candidatePath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	refPath, err := filepath.Abs(refSourceD1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve reference path: %v\n", err)
		os.Exit(1)
	}

	refOut, err := runProgram(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference execution failed: %v\n", err)
		os.Exit(1)
	}
	candOut, err := runProgram(candidatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate execution failed: %v\n", err)
		os.Exit(1)
	}

	if refOut != candOut {
		fmt.Fprintf(os.Stderr, "output mismatch:\nexpected:\n%s\nactual:\n%s\n", refOut, candOut)
		os.Exit(1)
	}

	fmt.Println("All tests passed")
}

func runProgram(path string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	var out bytes.Buffer
	var combined bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return strings.TrimSpace(out.String()), nil
}

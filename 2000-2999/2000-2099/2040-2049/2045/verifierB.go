package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "2000-2999/2000-2099/2040-2049/2045/2045B.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}
	if len(input) == 0 {
		fail("empty input")
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refAns, err := parseAnswer(refOut)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candAns, err := parseAnswer(candOut)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	if refAns != candAns {
		fail("outputs differ: expected %d got %d", refAns, candAns)
	}

	fmt.Println("OK")
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2045B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseAnswer(out string) (int64, error) {
	reader := strings.NewReader(out)
	var val int64
	if _, err := fmt.Fscan(reader, &val); err != nil {
		return 0, fmt.Errorf("failed to read integer: %v", err)
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return 0, fmt.Errorf("unexpected extra output: %s", extra)
	}
	return val, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

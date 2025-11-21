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

const refSource = "2000-2999/2000-2099/2050-2059/2050/2050E.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
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
	refAnswers, err := parseAnswers(refOut)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candAnswers, err := parseAnswers(candOut)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	if len(refAnswers) != len(candAnswers) {
		fail("expected %d answers, got %d", len(refAnswers), len(candAnswers))
	}
	for i := range refAnswers {
		if refAnswers[i] != candAnswers[i] {
			fail("test %d: expected %d got %d", i+1, refAnswers[i], candAnswers[i])
		}
	}

	fmt.Println("OK")
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2050E-ref-*")
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

func parseAnswers(out string) ([]int64, error) {
	reader := strings.NewReader(out)
	var res []int64
	for {
		var val int64
		_, err := fmt.Fscan(reader, &val)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to parse integer: %v", err)
		}
		res = append(res, val)
	}
	return res, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

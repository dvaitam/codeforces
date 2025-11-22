package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidatePath := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	testCount, err := parseTestCount(inputData)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, cleanupRef, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := prepareCandidate(candidatePath)
	if err != nil {
		fail("failed to prepare candidate: %v", err)
	}
	defer cleanupCand()

	refOut, refErr, err := runProgram(refBin, inputData)
	if err != nil {
		fail("reference runtime error: %v\n%s", err, refErr)
	}
	expected, err := parseOutputs(refOut, testCount)
	if err != nil {
		fail("failed to parse reference output: %v\noutput:\n%s", err, refOut)
	}

	candOut, candErr, err := runProgram(candBin, inputData)
	if err != nil {
		fail("candidate runtime error: %v\n%s", err, candErr)
	}
	got, err := parseOutputs(candOut, testCount)
	if err != nil {
		fail("invalid candidate output: %v\noutput:\n%s", err, candOut)
	}

	for i := 0; i < testCount; i++ {
		if got[i] != expected[i] {
			fail("test case %d: wrong answer, expected %d got %d", i+1, expected[i], got[i])
		}
	}

	fmt.Println("OK")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func parseTestCount(data []byte) (int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "2056D.go")

	tmp, err := os.CreateTemp("", "2056D-ref-*")
	if err != nil {
		return "", nil, err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", nil, fmt.Errorf("%v\n%s", err, out.String())
	}
	cleanup := func() {
		os.Remove(tmp.Name())
	}
	return tmp.Name(), cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", nil, err
		}
		tmp, err := os.CreateTemp("", "2056D-cand-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), abs)
		cmd.Dir = filepath.Dir(abs)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProgram(path string, input []byte) (string, string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.String(), errBuf.String(), err
}

func parseOutputs(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i := 0; i < t; i++ {
		v, err := strconv.ParseInt(tokens[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tokens[i])
		}
		res[i] = v
	}
	return res, nil
}

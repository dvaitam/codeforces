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

	n, err := readN(inputData)
	if err != nil {
		fail("failed to parse n from input: %v", err)
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
	expCost, expVerts, err := parseOutput(refOut, n)
	if err != nil {
		fail("failed to parse reference output: %v\noutput:\n%s", err, refOut)
	}

	candOut, candErr, err := runProgram(candBin, inputData)
	if err != nil {
		fail("candidate runtime error: %v\n%s", err, candErr)
	}
	gotCost, gotVerts, err := parseOutput(candOut, n)
	if err != nil {
		fail("invalid candidate output: %v\noutput:\n%s", err, candOut)
	}

	if gotCost != expCost {
		fail("wrong minimum cost: expected %d got %d", expCost, gotCost)
	}
	if len(gotVerts) != len(expVerts) {
		fail("wrong vertex count: expected %d got %d", len(expVerts), len(gotVerts))
	}
	for i, v := range expVerts {
		if gotVerts[i] != v {
			fail("vertices mismatch at position %d: expected %d got %d", i+1, v, gotVerts[i])
		}
	}

	fmt.Println("OK")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func readN(data []byte) (int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, err
	}
	return n, nil
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1120D.go")

	tmp, err := os.CreateTemp("", "1120D-ref-*")
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
		tmp, err := os.CreateTemp("", "1120D-cand-*")
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

func parseOutput(out string, n int) (int64, []int, error) {
	tokens := strings.Fields(out)
	if len(tokens) < 2 {
		return 0, nil, fmt.Errorf("expected at least 2 tokens, got %d", len(tokens))
	}
	cost, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return 0, nil, fmt.Errorf("invalid minimum cost %q", tokens[0])
	}
	k, err := strconv.Atoi(tokens[1])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid vertex count %q", tokens[1])
	}
	if k < 0 {
		return 0, nil, fmt.Errorf("negative vertex count %d", k)
	}
	if len(tokens) != 2+k {
		return 0, nil, fmt.Errorf("expected %d tokens, got %d", 2+k, len(tokens))
	}
	verts := make([]int, k)
	prev := 0
	for i := 0; i < k; i++ {
		v, err := strconv.Atoi(tokens[2+i])
		if err != nil {
			return 0, nil, fmt.Errorf("invalid vertex index %q", tokens[2+i])
		}
		if v < 1 || v > n {
			return 0, nil, fmt.Errorf("vertex %d out of range [1,%d]", v, n)
		}
		if i > 0 && v <= prev {
			return 0, nil, fmt.Errorf("vertices not in strictly increasing order")
		}
		prev = v
		verts[i] = v
	}
	return cost, verts, nil
}

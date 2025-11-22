package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSourceC = "2000-2999/2000-2099/2030-2039/2035/2035C.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}
	ns, err := parseN(input)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runCommand(exec.Command(refBin), input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refAns, _, err := parseOutput(refOut, ns)
	if err != nil {
		fail("failed to parse reference output: %v\n%s", err, refOut)
	}

	userOut, err := runCommand(commandFor(candidate), input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	userAns, userPerms, err := parseOutput(userOut, ns)
	if err != nil {
		fail("failed to parse candidate output: %v\n%s", err, userOut)
	}

	if len(userAns) != len(refAns) {
		fail("wrong number of test cases in output: expected %d got %d", len(refAns), len(userAns))
	}

	for i := range ns {
		if len(userPerms[i]) != ns[i] {
			fail("test %d: expected permutation of length %d, got %d numbers", i+1, ns[i], len(userPerms[i]))
		}
		if !isValidPermutation(userPerms[i]) {
			fail("test %d: permutation is invalid", i+1)
		}
		kVal := computeK(userPerms[i])
		if kVal != userAns[i] {
			fail("test %d: reported k=%d but permutation yields %d", i+1, userAns[i], kVal)
		}
		if userAns[i] != refAns[i] {
			fail("test %d: k not optimal, expected %d got %d", i+1, refAns[i], userAns[i])
		}
	}

	fmt.Println("OK")
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2035C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceC))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runCommand(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if msg == "" {
			msg = stdout.String()
		}
		return "", fmt.Errorf("%v\n%s", err, msg)
	}
	return stdout.String(), nil
}

func parseN(input []byte) ([]int, error) {
	fields := strings.Fields(string(input))
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	idx := 0
	t, err := strconv.Atoi(fields[idx])
	if err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	idx++
	ns := make([]int, t)
	for i := 0; i < t; i++ {
		if idx >= len(fields) {
			return nil, fmt.Errorf("unexpected end of input at test %d", i+1)
		}
		n, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("invalid n at test %d: %v", i+1, err)
		}
		ns[i] = n
		idx++
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("unexpected extra tokens in input")
	}
	return ns, nil
}

func parseOutput(out string, ns []int) ([]int64, [][]int, error) {
	fields := strings.Fields(out)
	ans := make([]int64, len(ns))
	perms := make([][]int, len(ns))
	idx := 0
	for i, n := range ns {
		if idx >= len(fields) {
			return nil, nil, fmt.Errorf("not enough numbers for test %d", i+1)
		}
		val, err := strconv.ParseInt(fields[idx], 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid answer for test %d: %v", i+1, err)
		}
		ans[i] = val
		idx++
		if idx+n > len(fields) {
			return nil, nil, fmt.Errorf("not enough permutation numbers for test %d", i+1)
		}
		perm := make([]int, n)
		for j := 0; j < n; j++ {
			val, err := strconv.Atoi(fields[idx+j])
			if err != nil {
				return nil, nil, fmt.Errorf("invalid permutation number at test %d position %d: %v", i+1, j+1, err)
			}
			perm[j] = val
		}
		perms[i] = perm
		idx += n
	}
	if idx != len(fields) {
		return nil, nil, fmt.Errorf("unexpected extra output data")
	}
	return ans, perms, nil
}

func isValidPermutation(p []int) bool {
	n := len(p)
	seen := make([]bool, n+1)
	for _, v := range p {
		if v < 1 || v > n || seen[v] {
			return false
		}
		seen[v] = true
	}
	return true
}

func computeK(p []int) int64 {
	var k int64
	for i, v := range p {
		val := int64(v)
		if (i+1)%2 == 1 {
			k &= val
		} else {
			k |= val
		}
	}
	return k
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

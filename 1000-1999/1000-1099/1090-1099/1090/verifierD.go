package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "1000-1999/1000-1099/1090-1099/1090/1090D.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	n, pairs, err := parseProblemInput(input)
	if err != nil {
		fail("failed to parse input: %v", err)
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
	refStatus, err := parseDecisionOnly(refOut)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	userOut, err := runProgram(candidate, input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	userStatus, first, second, err := parseCandidateOutput(userOut, n)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	if refStatus == "NO" {
		if userStatus != "NO" {
			fail("expected answer NO, but candidate produced YES")
		}
		fmt.Println("OK")
		return
	}

	if userStatus != "YES" {
		fail("expected answer YES, but candidate produced NO")
	}

	if err := validateArrays(n, pairs, first, second); err != nil {
		fail("invalid arrays: %v", err)
	}

	fmt.Println("OK")
}

func parseProblemInput(data []byte) (int, [][2]int, error) {
	reader := bytes.NewReader(data)
	in := bufio.NewReader(reader)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return 0, nil, err
	}
	pairs := make([][2]int, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Fscan(in, &pairs[i][0], &pairs[i][1]); err != nil {
			return 0, nil, err
		}
	}
	return n, pairs, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1090D-ref-*")
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

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseDecisionOnly(out string) (string, error) {
	reader := strings.NewReader(out)
	var status string
	if _, err := fmt.Fscan(reader, &status); err != nil {
		return "", fmt.Errorf("missing decision: %v", err)
	}
	status = strings.TrimSpace(status)
	if status != "YES" && status != "NO" {
		return "", fmt.Errorf("invalid decision %q", status)
	}
	return status, nil
}

func parseCandidateOutput(out string, n int) (string, []int, []int, error) {
	reader := strings.NewReader(out)
	var status string
	if _, err := fmt.Fscan(reader, &status); err != nil {
		return "", nil, nil, fmt.Errorf("missing decision: %v", err)
	}
	status = strings.TrimSpace(status)
	switch status {
	case "NO":
		var extra string
		if _, err := fmt.Fscan(reader, &extra); err == nil {
			return "", nil, nil, fmt.Errorf("unexpected extra output after NO")
		}
		return "NO", nil, nil, nil
	case "YES":
		first := make([]int, n)
		second := make([]int, n)
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(reader, &first[i]); err != nil {
				return "", nil, nil, fmt.Errorf("failed to read first array element %d: %v", i+1, err)
			}
		}
		for i := 0; i < n; i++ {
			if _, err := fmt.Fscan(reader, &second[i]); err != nil {
				return "", nil, nil, fmt.Errorf("failed to read second array element %d: %v", i+1, err)
			}
		}
		var extra string
		if _, err := fmt.Fscan(reader, &extra); err == nil {
			return "", nil, nil, fmt.Errorf("unexpected extra output after arrays")
		}
		return "YES", first, second, nil
	default:
		return "", nil, nil, fmt.Errorf("decision must be YES or NO, got %q", status)
	}
}

func validateArrays(n int, pairs [][2]int, first, second []int) error {
	if len(first) != n || len(second) != n {
		return fmt.Errorf("arrays must have length %d", n)
	}

	seenFirst := make([]bool, n+1)
	for i, val := range first {
		if val < 1 || val > n {
			return fmt.Errorf("first array value at position %d is %d, out of range [1,%d]", i+1, val, n)
		}
		if seenFirst[val] {
			return fmt.Errorf("first array values must be distinct, value %d repeats", val)
		}
		seenFirst[val] = true
	}

	seenSecond := make([]bool, n+1)
	hasDuplicate := false
	for i, val := range second {
		if val < 1 || val > n {
			return fmt.Errorf("second array value at position %d is %d, out of range [1,%d]", i+1, val, n)
		}
		if seenSecond[val] {
			hasDuplicate = true
		} else {
			seenSecond[val] = true
		}
	}
	if !hasDuplicate {
		return fmt.Errorf("second array must contain at least one repeated value")
	}

	for _, p := range pairs {
		a := p[0] - 1
		b := p[1] - 1
		if compare(first[a], first[b]) != compare(second[a], second[b]) {
			return fmt.Errorf("comparison mismatch for positions %d and %d", p[0], p[1])
		}
	}
	return nil
}

func compare(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

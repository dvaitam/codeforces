package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const maxVal = 500000 + 5

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	m, k, n, a, need, unique, err := parseInput(inputData)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(refBin, inputData)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refPossible, err := parseReferenceOutcome(refOut)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	userOut, err := runProgram(candidate, inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	impossible, removed, err := parseCandidateOutput(userOut)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	if impossible {
		if refPossible {
			fail("candidate claims impossible but solution exists")
		}
		fmt.Println("OK")
		return
	}

	if !refPossible {
		fail("candidate claims possible but reference says impossible")
	}

	if err := validateSolution(a, need, unique, m, k, n, removed); err != nil {
		fail("invalid solution: %v", err)
	}

	fmt.Println("OK")
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func parseInput(data []byte) (int, int, int, []int, []int, []int, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var m, k, n, s int
	if _, err := fmt.Fscan(reader, &m, &k, &n, &s); err != nil {
		return 0, 0, 0, nil, nil, nil, err
	}
	a := make([]int, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Fscan(reader, &a[i]); err != nil {
			return 0, 0, 0, nil, nil, nil, err
		}
	}
	need := make([]int, maxVal)
	var unique []int
	for i := 0; i < s; i++ {
		var b int
		if _, err := fmt.Fscan(reader, &b); err != nil {
			return 0, 0, 0, nil, nil, nil, err
		}
		if need[b] == 0 {
			unique = append(unique, b)
		}
		need[b]++
	}
	return m, k, n, a, need, unique, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1120A-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean("1000-1999/1100-1199/1120-1129/1120/1120A.go"))
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
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseReferenceOutcome(out string) (bool, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var token string
	if _, err := fmt.Fscan(reader, &token); err != nil {
		return false, fmt.Errorf("reference produced no output")
	}
	if token == "-1" {
		return false, nil
	}
	return true, nil
}

func parseCandidateOutput(out string) (bool, []int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var first string
	if _, err := fmt.Fscan(reader, &first); err != nil {
		return false, nil, fmt.Errorf("failed to read first token: %v", err)
	}
	if first == "-1" {
		var extra string
		if _, err := fmt.Fscan(reader, &extra); err == nil {
			return false, nil, fmt.Errorf("extra output after -1")
		}
		return true, nil, nil
	}
	d, err := strconv.Atoi(first)
	// if candidate prints d as string, convert
	if err != nil {
		return false, nil, fmt.Errorf("invalid removal count %q", first)
	}
	if d < 0 {
		return false, nil, fmt.Errorf("negative removal count")
	}
	removed := make([]int, d)
	for i := 0; i < d; i++ {
		if _, err := fmt.Fscan(reader, &removed[i]); err != nil {
			return false, nil, fmt.Errorf("failed to read position %d: %v", i+1, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return false, nil, fmt.Errorf("extra output detected")
	}
	return false, removed, nil
}

func validateSolution(a []int, need []int, unique []int, m, k, n int, removed []int) error {
	if len(removed) > m {
		return fmt.Errorf("too many removals")
	}
	removedMark := make([]bool, m)
	for _, pos := range removed {
		if pos < 1 || pos > m {
			return fmt.Errorf("position %d out of range", pos)
		}
		if removedMark[pos-1] {
			return fmt.Errorf("duplicate position %d", pos)
		}
		removedMark[pos-1] = true
	}
	remaining := m - len(removed)
	if remaining < n*k {
		return fmt.Errorf("not enough flowers remaining")
	}

	filtered := make([]int, 0, remaining)
	for i := 0; i < m; i++ {
		if !removedMark[i] {
			filtered = append(filtered, a[i])
		}
	}
	blocks := len(filtered) / k
	if blocks < n {
		return fmt.Errorf("not enough workpieces after removal")
	}

	freq := make([]int, maxVal)
	visited := make([]int, 0, len(unique))
	for block := 0; block < blocks; block++ {
		visited = visited[:0]
		start := block * k
		for i := start; i < start+k; i++ {
			val := filtered[i]
			if need[val] == 0 {
				continue
			}
			if freq[val] == 0 {
				visited = append(visited, val)
			}
			freq[val]++
		}
		ok := true
		for _, typ := range unique {
			if freq[typ] < need[typ] {
				ok = false
				break
			}
		}
		for _, typ := range visited {
			freq[typ] = 0
		}
		if ok {
			return nil
		}
	}
	return fmt.Errorf("no valid workpiece found")
}

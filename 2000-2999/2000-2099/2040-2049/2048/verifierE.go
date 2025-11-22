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

const refSource = "2000-2999/2000-2099/2040-2049/2048/2048E.go"

type testCase struct {
	n int
	m int
}

func main() {
	if len(os.Args) != 2 {
		fail("usage: go run verifierE.go /path/to/candidate")
	}
	candidate := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	tests, err := parseInput(inputData)
	if err != nil {
		fail("failed to parse input: %v", err)
	}

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(exec.Command(refBin), inputData)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refVerdicts, err := parseOutput(refOut, tests, false)
	if err != nil {
		fail("invalid reference output: %v", err)
	}

	candOut, err := runProgram(commandFor(candidate), inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candVerdicts, err := parseOutput(candOut, tests, true)
	if err != nil {
		fail("invalid candidate output: %v", err)
	}

	for i := range tests {
		if refVerdicts[i] == "yes" {
			if candVerdicts[i] != "yes" {
				fail("test %d: reference found a solution but candidate printed NO", i+1)
			}
		} else {
			if candVerdicts[i] == "yes" {
				fail("test %d: reference printed NO but candidate printed YES", i+1)
			}
		}
	}

	fmt.Println("OK")
}

func parseInput(data []byte) ([]testCase, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &tests[i].n, &tests[i].m); err != nil {
			return nil, err
		}
	}
	return tests, nil
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2048E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, combined.String())
	}
	return tmp.Name(), nil
}

func runProgram(cmd *exec.Cmd, input []byte) (string, error) {
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func parseOutput(out string, tests []testCase, validate bool) ([]string, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	verdicts := make([]string, len(tests))
	for idx, tc := range tests {
		token, err := readToken(reader)
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("not enough outputs for test %d", idx+1)
			}
			return nil, err
		}
		ver := strings.ToLower(token)
		if ver != "yes" && ver != "no" {
			return nil, fmt.Errorf("test %d: expected YES/NO, got %q", idx+1, token)
		}
		verdicts[idx] = ver
		if ver == "no" {
			continue
		}
		if err := readAndValidateMatrix(reader, tc.n, tc.m, validate); err != nil {
			return nil, fmt.Errorf("test %d: %v", idx+1, err)
		}
	}
	if extra, err := readToken(reader); err == nil {
		return nil, fmt.Errorf("unexpected extra output token %q", extra)
	} else if err != io.EOF {
		return nil, err
	}
	return verdicts, nil
}

func readAndValidateMatrix(r *bufio.Reader, n, m int, validate bool) error {
	totalRows := 2 * n
	matrix := make([][]int, totalRows)
	for i := 0; i < totalRows; i++ {
		row := make([]int, m)
		for j := 0; j < m; j++ {
			tok, err := readToken(r)
			if err != nil {
				if err == io.EOF {
					return fmt.Errorf("expected %d rows of %d numbers", totalRows, m)
				}
				return err
			}
			val, convErr := strconv.Atoi(tok)
			if convErr != nil {
				return fmt.Errorf("invalid integer %q", tok)
			}
			row[j] = val
		}
		matrix[i] = row
	}
	if !validate {
		return nil
	}
	return validateSolution(matrix, n, m)
}

func validateSolution(matrix [][]int, n, m int) error {
	// Verify color range and detect monochromatic cycle per color via DSU.
	colors := make(map[int]*dsu)
	leftCount := 2 * n
	totalNodes := leftCount + m

	for i := 0; i < leftCount; i++ {
		row := matrix[i]
		if len(row) != m {
			return fmt.Errorf("row %d length mismatch", i+1)
		}
		for j := 0; j < m; j++ {
			c := row[j]
			if c < 1 || c > n {
				return fmt.Errorf("color out of range at row %d col %d: %d", i+1, j+1, c)
			}
			d := colors[c]
			if d == nil {
				d = newDSU(totalNodes)
				colors[c] = d
			}
			u := i
			v := leftCount + j
			if !d.union(u, v) {
				return fmt.Errorf("monochromatic cycle detected with color %d", c)
			}
		}
	}
	return nil
}

type dsu struct {
	p []int
	r []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &dsu{p: p, r: r}
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

func (d *dsu) union(a, b int) bool {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return false
	}
	if d.r[a] < d.r[b] {
		a, b = b, a
	}
	d.p[b] = a
	if d.r[a] == d.r[b] {
		d.r[a]++
	}
	return true
}

func readToken(r *bufio.Reader) (string, error) {
	var b strings.Builder
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if ch > ' ' {
			b.WriteByte(ch)
			break
		}
	}
	for {
		ch, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				return b.String(), nil
			}
			return "", err
		}
		if ch <= ' ' {
			return b.String(), nil
		}
		b.WriteByte(ch)
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

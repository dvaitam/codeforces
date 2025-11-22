package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2000-2099/2020-2029/2029/2029I.go"

type testCase struct {
	n int
	m int
	k int64
	a []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := renderInput(tests)

	refOut, err := runWithInput(exec.Command(refBin), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	candCmd := commandFor(candidate)
	candOut, err := runWithInput(candCmd, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	expect, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		expLine := expect[i]
		gotLine := got[i]
		if len(expLine) != len(gotLine) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d numbers, got %d\ninput:\n%s", i+1, len(expLine), len(gotLine), formatSingleInput(tests[i]))
			os.Exit(1)
		}
		for j := range expLine {
			if expLine[j] != gotLine[j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d position %d: expected %d got %d\ninput:\n%s", i+1, j+1, expLine[j], gotLine[j], formatSingleInput(tests[i]))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2029I-ref-*")
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

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, tests []testCase) ([][]int64, error) {
	fields := strings.Fields(output)
	totalNeed := 0
	for _, tc := range tests {
		totalNeed += tc.m
	}
	if len(fields) < totalNeed {
		return nil, fmt.Errorf("expected %d numbers, got %d", totalNeed, len(fields))
	}
	if len(fields) > totalNeed {
		return nil, fmt.Errorf("extra output detected after %d numbers", totalNeed)
	}

	res := make([][]int64, len(tests))
	ptr := 0
	for i, tc := range tests {
		line := make([]int64, tc.m)
		for j := 0; j < tc.m; j++ {
			val, err := strconv.ParseInt(fields[ptr], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse number %d (%q): %v", ptr+1, fields[ptr], err)
			}
			line[j] = val
			ptr++
		}
		res[i] = line
	}
	return res, nil
}

func renderInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func formatSingleInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 3, m: 2, k: 1, a: []int64{1, 2, 2}},
		{n: 3, m: 2, k: 2, a: []int64{3, 2, 2}},
		{n: 1, m: 2, k: 2, a: []int64{1}},
		{n: 10, m: 2, k: 1, a: []int64{10, 1, 1, 1, 1, 10, 1, 1, 1, 1}},
		{n: 6, m: 8, k: 2, a: []int64{1, 1, 4, 5, 1, 3}},
	}

	totalProd := 0
	for _, tc := range tests {
		totalProd += tc.n * tc.m
	}

	const maxTotal = 20000
	rng := rand.New(rand.NewSource(2029))
	for totalProd < maxTotal {
		n := rng.Intn(80) + 1
		maxM := maxTotal / n
		if maxM > 200 {
			maxM = 200
		}
		m := rng.Intn(maxM) + 1
		if totalProd+n*m > maxTotal {
			m = (maxTotal - totalProd) / n
			if m == 0 {
				break
			}
		}
		k := int64(rng.Intn(100000) + 1)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = int64(rng.Intn(100000) + 1)
		}
		tests = append(tests, testCase{n: n, m: m, k: k, a: a})
		totalProd += n * m
		if len(tests) > 150 {
			break
		}
	}

	// Add one large-n small-m case near limits if space permits.
	if totalProd < maxTotal {
		n := 5000
		m := (maxTotal - totalProd) / n
		if m > 0 {
			if m > 4 {
				m = 4
			}
			a := make([]int64, n)
			for i := 0; i < n; i++ {
				a[i] = int64((i % 7) + 1)
			}
			tests = append(tests, testCase{n: n, m: m, k: 3, a: a})
		}
	}

	return tests
}

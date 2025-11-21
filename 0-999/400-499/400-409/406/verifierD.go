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
	"time"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	outPath := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, "406D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(out))
	}
	return outPath, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(x []int64, y []int64, queries [][2]int) string {
	var sb strings.Builder
	n := len(x)
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", x[i], y[i])
	}
	fmt.Fprintf(&sb, "%d\n", len(queries))
	for _, q := range queries {
		fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int, expected)
	for i := 0; i < expected; i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", fields[i])
		}
		res[i] = v
	}
	return res, nil
}

func deterministicTests() []string {
	var tests []string
	// single hill
	tests = append(tests, buildInput(
		[]int64{1},
		[]int64{10},
		[][2]int{{1, 1}},
	))
	// two hills, mutual visibility
	tests = append(tests, buildInput(
		[]int64{1, 2},
		[]int64{5, 7},
		[][2]int{{1, 2}, {2, 1}},
	))
	// three hills monotone heights
	tests = append(tests, buildInput(
		[]int64{1, 3, 5},
		[]int64{3, 2, 4},
		[][2]int{{1, 2}, {1, 3}, {2, 3}},
	))
	// plateau-like
	tests = append(tests, buildInput(
		[]int64{1, 4, 7, 11},
		[]int64{10, 1, 10, 1},
		[][2]int{{1, 4}, {2, 3}, {1, 2}, {3, 4}},
	))
	// increasing x gaps
	tests = append(tests, buildInput(
		[]int64{1, 10, 20, 35, 60},
		[]int64{5, 15, 5, 25, 10},
		[][2]int{{1, 3}, {2, 4}, {4, 5}, {5, 1}},
	))
	return tests
}

func randomHills(n int, rnd *rand.Rand) ([]int64, []int64) {
	x := make([]int64, n)
	y := make([]int64, n)
	curX := int64(0)
	for i := 0; i < n; i++ {
		curX += int64(rnd.Intn(10) + 1)
		y[i] = rnd.Int63n(1_000_000_000_000) + 1
		x[i] = curX
	}
	return x, y
}

func randomQueries(n, m int, rnd *rand.Rand) [][2]int {
	qs := make([][2]int, m)
	for i := 0; i < m; i++ {
		a := rnd.Intn(n) + 1
		b := rnd.Intn(n) + 1
		qs[i] = [2]int{a, b}
	}
	return qs
}

func randomTests(count int) []string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 0, count)
	for i := 0; i < count; i++ {
		var n int
		switch {
		case i%25 == 0:
			n = rnd.Intn(100000) + 1
		case i%5 == 0:
			n = rnd.Intn(5000) + 1
		default:
			n = rnd.Intn(200) + 1
		}
		m := rnd.Intn(min(n, 2000)) + 1
		if i%20 == 0 {
			m = rnd.Intn(100000) + 1
		}
		x, y := randomHills(n, rnd)
		queries := randomQueries(n, m, rnd)
		tests = append(tests, buildInput(x, y, queries))
	}
	// add one maximal-ish case
	n := 100000
	m := 100000
	x := make([]int64, n)
	y := make([]int64, n)
	for i := 0; i < n; i++ {
		x[i] = int64(i + 1)
		y[i] = int64((i*i)%1_000_000 + 1)
	}
	rnd2 := rand.New(rand.NewSource(time.Now().UnixNano() + 1))
	queries := randomQueries(n, m, rnd2)
	tests = append(tests, buildInput(x, y, queries))
	return tests
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(200)...)

	for idx, input := range tests {
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		// parse header to get number of queries
		scanner := strings.NewReader(input)
		var nVal int
		fmt.Fscan(scanner, &nVal)
		for i := 0; i < nVal; i++ {
			var xi, yi int64
			fmt.Fscan(scanner, &xi, &yi)
		}
		var m int
		fmt.Fscan(scanner, &m)

		expVals, err := parseOutput(expOut, m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		for i := 0; i < m; i++ {
			if expVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "case %d query %d mismatch: expected %d got %d\n", idx+1, i+1, expVals[i], gotVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

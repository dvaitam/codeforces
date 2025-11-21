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
	cmd := exec.Command("go", "build", "-o", outPath, "516D.go")
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

func buildInput(n int, edges [][3]int64, queries []int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e[0], e[1], e[2])
	}
	fmt.Fprintf(&sb, "%d\n", len(queries))
	for i, l := range queries {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", l)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string, expected int) ([]int64, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(lines))
	}
	res := make([]int64, expected)
	for i := 0; i < expected; i++ {
		val, err := strconv.ParseInt(lines[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", lines[i])
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []string {
	var tests []string
	// smallest tree (line)
	tests = append(tests, buildInput(
		2,
		[][3]int64{{1, 2, 1}},
		[]int64{1, 2, 10},
	))
	// simple star
	tests = append(tests, buildInput(
		4,
		[][3]int64{{1, 2, 1}, {1, 3, 2}, {1, 4, 3}},
		[]int64{0, 1, 5, 10},
	))
	// chain with varying weights
	tests = append(tests, buildInput(
		5,
		[][3]int64{{1, 2, 1}, {2, 3, 2}, {3, 4, 3}, {4, 5, 4}},
		[]int64{1, 3, 6, 10},
	))
	// balanced-ish tree
	tests = append(tests, buildInput(
		7,
		[][3]int64{{1, 2, 1}, {1, 3, 1}, {2, 4, 2}, {2, 5, 2}, {3, 6, 2}, {3, 7, 2}},
		[]int64{0, 2, 4, 6},
	))
	return tests
}

func randomTree(n int, rnd *rand.Rand) [][3]int64 {
	edges := make([][3]int64, 0, n-1)
	for i := 2; i <= n; i++ {
		parent := rnd.Intn(i-1) + 1
		weight := rnd.Int63n(1_000_000) + 1
		edges = append(edges, [3]int64{int64(parent), int64(i), weight})
	}
	return edges
}

func randomQueries(rnd *rand.Rand, q int) []int64 {
	res := make([]int64, q)
	for i := 0; i < q; i++ {
		res[i] = rnd.Int63n(1_000_000_000_000) + 1
	}
	return res
}

func randomTests(count int) []string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 0, count)
	for i := 0; i < count; i++ {
		var n int
		switch {
		case i%30 == 0:
			n = rnd.Intn(100000) + 2
		case i%5 == 0:
			n = rnd.Intn(5000) + 2
		default:
			n = rnd.Intn(200) + 2
		}
		edges := randomTree(n, rnd)
		q := rnd.Intn(50) + 1
		queries := randomQueries(rnd, q)
		tests = append(tests, buildInput(n, edges, queries))
	}
	return tests
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
		reader := strings.NewReader(input)
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to read n: %v\n", idx+1, err)
			os.Exit(1)
		}
		for i := 0; i < n-1; i++ {
			var x, y int
			var w int64
			fmt.Fscan(reader, &x, &y, &w)
		}
		var q int
		if _, err := fmt.Fscan(reader, &q); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to read q: %v\n", idx+1, err)
			os.Exit(1)
		}
		for i := 0; i < q; i++ {
			var tmp int64
			fmt.Fscan(reader, &tmp)
		}

		expVals, err := parseOutput(expOut, q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		for i := 0; i < q; i++ {
			if expVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "case %d, query %d mismatch: expected %d got %d\n", idx+1, i+1, expVals[i], gotVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

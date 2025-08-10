package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// runBin executes the provided binary or go source file with the given input
// and returns its stdout as a trimmed string.
func runBin(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// genTest generates a random valid test for problem 1517B.
func genTest(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := rng.Intn(4) + 1
		m := rng.Intn(4) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprint(rng.Intn(10) + 1))
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

// parseInput converts the textual input into structured test cases.
type testCase struct {
	n, m int
	rows [][]int
}

func parseInput(inp string) ([]testCase, error) {
	fields := strings.Fields(inp)
	idx := 0
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	t, err := strconv.Atoi(fields[idx])
	if err != nil {
		return nil, err
	}
	idx++
	cases := make([]testCase, t)
	for c := 0; c < t; c++ {
		if idx+1 >= len(fields) {
			return nil, fmt.Errorf("unexpected end of input")
		}
		n, _ := strconv.Atoi(fields[idx])
		m, _ := strconv.Atoi(fields[idx+1])
		idx += 2
		rows := make([][]int, n)
		for i := 0; i < n; i++ {
			rows[i] = make([]int, m)
			for j := 0; j < m; j++ {
				if idx >= len(fields) {
					return nil, fmt.Errorf("unexpected end of input")
				}
				v, err := strconv.Atoi(fields[idx])
				if err != nil {
					return nil, err
				}
				rows[i][j] = v
				idx++
			}
		}
		cases[c] = testCase{n: n, m: m, rows: rows}
	}
	return cases, nil
}

// parseOutput parses the candidate output according to the provided test cases.
func parseOutput(out string, cases []testCase) ([][][]int, error) {
	fields := strings.Fields(out)
	idx := 0
	res := make([][][]int, len(cases))
	for c, tc := range cases {
		matrix := make([][]int, tc.n)
		for i := 0; i < tc.n; i++ {
			row := make([]int, tc.m)
			for j := 0; j < tc.m; j++ {
				if idx >= len(fields) {
					return nil, fmt.Errorf("not enough output data")
				}
				v, err := strconv.Atoi(fields[idx])
				if err != nil {
					return nil, fmt.Errorf("invalid integer: %v", fields[idx])
				}
				row[j] = v
				idx++
			}
			matrix[i] = row
		}
		res[c] = matrix
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("too much output data")
	}
	return res, nil
}

// checkCase validates a single test case output.
func checkCase(tc testCase, out [][]int) bool {
	// verify each row is a permutation of the input row
	for i := 0; i < tc.n; i++ {
		a := append([]int(nil), tc.rows[i]...)
		b := append([]int(nil), out[i]...)
		sort.Ints(a)
		sort.Ints(b)
		for j := 0; j < tc.m; j++ {
			if a[j] != b[j] {
				return false
			}
		}
	}
	// compute sum of column minima in candidate output
	sumCand := 0
	for j := 0; j < tc.m; j++ {
		minVal := out[0][j]
		for i := 1; i < tc.n; i++ {
			if out[i][j] < minVal {
				minVal = out[i][j]
			}
		}
		sumCand += minVal
	}
	// compute expected minimal sum
	flat := make([]int, 0, tc.n*tc.m)
	for i := 0; i < tc.n; i++ {
		flat = append(flat, tc.rows[i]...)
	}
	sort.Ints(flat)
	sumExp := 0
	for i := 0; i < tc.m; i++ {
		sumExp += flat[i]
	}
	return sumCand == sumExp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genTest(rng)
		cases, err := parseInput(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse generated test: %v\n", err)
			os.Exit(1)
		}
		got, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n" + input)
			os.Exit(1)
		}
		outs, err := parseOutput(got, cases)
		if err != nil {
			fmt.Printf("failed to parse output on test %d: %v\n", i+1, err)
			fmt.Println("input:\n" + input)
			fmt.Println("output:\n" + got)
			os.Exit(1)
		}
		for caseIdx, tc := range cases {
			if !checkCase(tc, outs[caseIdx]) {
				fmt.Printf("wrong answer on test %d\n", i+1)
				fmt.Println("input:\n" + input)
				fmt.Println("output:\n" + got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}

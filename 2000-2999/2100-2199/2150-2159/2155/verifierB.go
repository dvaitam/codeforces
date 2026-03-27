package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// solveReference is the correct solver for 2155B, embedded directly.
func solveReference(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var out bytes.Buffer
	writer := bufio.NewWriter(&out)

	getCell := func(m, n int) (int, int) {
		r := m / n
		c := m % n
		if r%2 == 1 {
			c = n - 1 - c
		}
		return r, c
	}

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return ""
	}

	for tc := 0; tc < t; tc++ {
		var n, k int
		fmt.Fscan(reader, &n, &k)

		C := n*n - k
		if C == 1 {
			fmt.Fprintln(writer, "NO")
			continue
		}

		fmt.Fprintln(writer, "YES")
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]byte, n)
			for j := 0; j < n; j++ {
				if i%2 == 0 {
					if j == n-1 {
						grid[i][j] = 'D'
					} else {
						grid[i][j] = 'R'
					}
				} else {
					if j == 0 {
						grid[i][j] = 'D'
					} else {
						grid[i][j] = 'L'
					}
				}
			}
		}

		if C >= 2 {
			r1, c1 := getCell(C-2, n)
			r2, c2 := getCell(C-1, n)
			if r1 == r2 {
				if c1 < c2 {
					grid[r2][c2] = 'L'
				} else {
					grid[r2][c2] = 'R'
				}
			} else {
				if r1 < r2 {
					grid[r2][c2] = 'U'
				} else {
					grid[r2][c2] = 'D'
				}
			}
		}

		for i := 0; i < n; i++ {
			fmt.Fprintln(writer, string(grid[i]))
		}
	}
	writer.Flush()
	return out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()
	for i, input := range tests {
		wantOut := solveReference(input)
		want, err := parseOutput(input, wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := parseOutput(input, gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, input, gotOut)
			os.Exit(1)
		}

		if err := compareSolutions(want, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s", i+1, err, input, wantOut, gotOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

type testCaseInput struct {
	n, k int
}

type outputCase struct {
	ok   bool
	grid [][]byte
}

type outputData struct {
	cases []outputCase
}

func parseInput(input string) ([]testCaseInput, error) {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	nextInt := func() (int, error) {
		if !sc.Scan() {
			return 0, fmt.Errorf("unexpected EOF")
		}
		return strconv.Atoi(sc.Text())
	}
	t, err := nextInt()
	if err != nil {
		return nil, err
	}
	cases := make([]testCaseInput, t)
	for i := 0; i < t; i++ {
		n, err := nextInt()
		if err != nil {
			return nil, err
		}
		k, err := nextInt()
		if err != nil {
			return nil, err
		}
		cases[i] = testCaseInput{n: n, k: k}
	}
	return cases, nil
}

func parseOutput(input string, out string) (outputData, error) {
	cases, err := parseInput(input)
	if err != nil {
		return outputData{}, err
	}
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanLines)
	idx := 0
	result := make([]outputCase, len(cases))
	for idx < len(cases) {
		if !sc.Scan() {
			return outputData{}, fmt.Errorf("case %d: missing verdict line", idx+1)
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		lineUpper := strings.ToUpper(line)
		if lineUpper == "NO" {
			result[idx] = outputCase{ok: false}
			idx++
			continue
		}
		if lineUpper != "YES" {
			return outputData{}, fmt.Errorf("case %d: expected YES/NO, got %q", idx+1, line)
		}
		grid := make([][]byte, cases[idx].n)
		for i := 0; i < cases[idx].n; i++ {
			if !sc.Scan() {
				return outputData{}, fmt.Errorf("case %d: incomplete grid", idx+1)
			}
			row := strings.TrimSpace(sc.Text())
			rowUpper := strings.ToUpper(row)
			if len(rowUpper) != cases[idx].n {
				return outputData{}, fmt.Errorf("case %d: row %d length mismatch", idx+1, i+1)
			}
			grid[i] = []byte(rowUpper)
		}
		result[idx] = outputCase{ok: true, grid: grid}
		idx++
	}
	return outputData{cases: result}, nil
}

func compareSolutions(want outputData, got outputData) error {
	if len(want.cases) != len(got.cases) {
		return fmt.Errorf("number of cases mismatch")
	}
	for i := range want.cases {
		if want.cases[i].ok != got.cases[i].ok {
			return fmt.Errorf("case %d: expected verdict %v got %v", i+1, want.cases[i].ok, got.cases[i].ok)
		}
		if want.cases[i].ok {
			if err := compareGrid(want.cases[i].grid, got.cases[i].grid); err != nil {
				return fmt.Errorf("case %d: %w", i+1, err)
			}
		}
	}
	return nil
}

func compareGrid(ref [][]byte, cand [][]byte) error {
	if len(ref) != len(cand) {
		return fmt.Errorf("grid size mismatch")
	}
	for i := range ref {
		if len(ref[i]) != len(cand[i]) {
			return fmt.Errorf("row %d length mismatch", i+1)
		}
		for j := range ref[i] {
			if ref[i][j] != cand[i][j] {
				return fmt.Errorf("cell (%d,%d) mismatch", i+1, j+1)
			}
		}
	}
	return nil
}

func generateTests() []string {
	var tests []string
	tests = append(tests, buildTest([]testCaseInput{
		{n: 2, k: 0},
		{n: 2, k: 3},
	}))

	tests = append(tests, buildTest([]testCaseInput{
		{n: 3, k: 5},
		{n: 3, k: 9},
	}))

	tests = append(tests, buildTest([]testCaseInput{
		{n: 4, k: 7},
		{n: 4, k: 13},
	}))

	tests = append(tests, buildTest([]testCaseInput{
		{n: 5, k: 18},
	}))

	tests = append(tests, buildTest([]testCaseInput{
		{n: 7, k: 20},
		{n: 7, k: 40},
		{n: 7, k: 48},
	}))

	tests = append(tests, buildTest([]testCaseInput{
		{n: 10, k: 50},
	}))

	tests = append(tests, buildTest([]testCaseInput{
		{n: 15, k: 120},
		{n: 15, k: 200},
	}))

	tests = append(tests, buildTest([]testCaseInput{
		{n: 20, k: 250},
	}))

	tests = append(tests, buildTest([]testCaseInput{
		{n: 25, k: 300},
		{n: 25, k: 500},
	}))

	tests = append(tests, buildTest([]testCaseInput{
		{n: 30, k: 450},
	}))

	tests = append(tests, buildTest([]testCaseInput{
		randomCase(50, 600),
		randomCase(80, 2000),
		randomCase(100, 9999),
	}))

	return tests
}

func buildTest(cases []testCaseInput) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, cs.k))
	}
	return sb.String()
}

func randomCase(n, k int) testCaseInput {
	if n < 2 {
		n = 2
	}
	if k < 0 {
		k = 0
	}
	if k > n*n {
		k = n * n
	}
	return testCaseInput{n: n, k: k}
}

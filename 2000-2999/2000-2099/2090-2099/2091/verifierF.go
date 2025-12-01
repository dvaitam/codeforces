package main

import (
	"bufio"
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

const refSource = "./2091F.go"

type testCase struct {
	n, m, d int
	grid    []string
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2091F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 3, m: 4, d: 1,
			grid: []string{
				"XX#X",
				"#XX#",
				"##X#",
			},
		},
		{
			n: 3, m: 4, d: 2,
			grid: []string{
				"XX#X",
				"#XX#",
				"##X#",
			},
		},
		{
			n: 3, m: 1, d: 3,
			grid: []string{
				"X",
				"X",
				"#",
			},
		},
		{
			n: 2, m: 2, d: 1,
			grid: []string{
				"##",
				"##",
			},
		},
		{
			n: 2, m: 3, d: 2,
			grid: []string{
				"###",
				"XXX",
			},
		},
	}
}

func randomGrid(rng *rand.Rand, n, m int, density float64) []string {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Float64() < density {
				row[j] = 'X'
			} else {
				row[j] = '#'
			}
		}
		grid[i] = string(row)
	}
	return grid
}

func randomTest(rng *rand.Rand, n, m int) testCase {
	d := rng.Intn(10) + 1
	if rng.Intn(5) == 0 {
		d = rng.Intn(50) + 1
	}
	density := 0.4 + rng.Float64()*0.4 // between 0.4 and 0.8
	grid := randomGrid(rng, n, m, density)
	return testCase{n: n, m: m, d: d, grid: grid}
}

func chainTest(n int, d int) testCase {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, 1)
		if i == n-1 || i == 0 {
			row[0] = 'X'
		} else {
			row[0] = 'X'
		}
		grid[i] = string(row)
	}
	return testCase{n: n, m: 1, d: d, grid: grid}
}

func bottomMissingTest(n, m, d int) testCase {
	grid := randomGrid(rand.New(rand.NewSource(int64(n*m+d))), n, m, 0.5)
	grid[n-1] = strings.Repeat("#", m)
	return testCase{n: n, m: m, d: d, grid: grid}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.d))
		for _, row := range tc.grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]string, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	res := make([]string, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("missing output for test %d", i+1)
		}
		res[i] = sc.Text()
	}
	if sc.Scan() {
		return nil, fmt.Errorf("extra output detected after %d testcases", t)
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := deterministicTests()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cellBudget := 300000 // keep runtime reasonable

	for len(tests) < 80 && cellBudget > 0 {
		n := rng.Intn(30) + 2
		m := rng.Intn(30) + 1
		if n*m > cellBudget {
			break
		}
		tests = append(tests, randomTest(rng, n, m))
		cellBudget -= n * m
	}

	// Some larger structured cases
	if cellBudget > 50000 {
		tests = append(tests, chainTest(2000, 1))
		cellBudget -= 2000
	}
	if cellBudget > 40000 {
		tests = append(tests, randomTest(rng, 200, 200))
		cellBudget -= 200 * 200
	}
	if cellBudget > 5000 {
		tests = append(tests, bottomMissingTest(50, 50, 5))
		cellBudget -= 50 * 50
	}

	input := buildInput(tests)

	wantOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	want, err := parseOutput(wantOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}

	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutput(gotOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if want[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\nn=%d m=%d d=%d\n", i+1, want[i], got[i], tests[i].n, tests[i].m, tests[i].d)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

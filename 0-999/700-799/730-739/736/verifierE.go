package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

const (
	referenceSolutionRel = "0-999/700-799/730-739/736/736E.go"
)

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "736E.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	m, n int
	top  []int
}

func randomValidCase(rng *rand.Rand, maxM int) testCase {
	m := rng.Intn(maxM) + 1
	points := make([]int, m)
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			switch rng.Intn(3) {
			case 0:
				points[i] += 2
			case 1:
				points[j] += 2
			default:
				points[i]++
				points[j]++
			}
		}
	}
	sorted := append([]int(nil), points...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] > sorted[j] })
	n := rng.Intn(m) + 1
	top := append([]int(nil), sorted[:n]...)
	return testCase{m: m, n: n, top: top}
}

func randomGeneralCase(rng *rand.Rand, maxM int) testCase {
	m := rng.Intn(maxM) + 1
	n := rng.Intn(m) + 1
	top := make([]int, n)
	maxPoints := 2 * (m - 1)
	if maxPoints < 0 {
		maxPoints = 0
	}
	for i := range top {
		top[i] = rng.Intn(maxPoints + 1)
	}
	sort.Slice(top, func(i, j int) bool { return top[i] > top[j] })
	return testCase{m: m, n: n, top: top}
}

func inputString(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.m, tc.n))
	for i, v := range tc.top {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "736E-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_736E")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseOutput(out string, m int) (string, []string, error) {
	reader := strings.NewReader(out)
	var verdict string
	if _, err := fmt.Fscan(reader, &verdict); err != nil {
		return "", nil, fmt.Errorf("failed to read verdict: %v\nfull output:\n%s", err, out)
	}
	verdict = strings.ToLower(verdict)
	if verdict == "no" {
		return "no", nil, nil
	}
	if verdict != "yes" {
		return "", nil, fmt.Errorf("expected 'yes' or 'no', got %q", verdict)
	}
	rows := make([]string, 0, m)
	for len(rows) < m {
		var row string
		if _, err := fmt.Fscan(reader, &row); err != nil {
			return "", nil, fmt.Errorf("expected %d rows after 'yes', got %d (%v)", m, len(rows), err)
		}
		rows = append(rows, row)
	}
	return "yes", rows, nil
}

func checkTable(tc testCase, rows []string) error {
	if len(rows) != tc.m {
		return fmt.Errorf("expected %d rows, got %d", tc.m, len(rows))
	}
	points := make([]int, tc.m)
	matrix := make([][]byte, tc.m)
	for i := 0; i < tc.m; i++ {
		row := rows[i]
		if len(row) != tc.m {
			return fmt.Errorf("row %d length mismatch: got %d want %d", i+1, len(row), tc.m)
		}
		matrix[i] = []byte(row)
	}
	for i := 0; i < tc.m; i++ {
		for j := 0; j < tc.m; j++ {
			ch := matrix[i][j]
			if 'a' <= ch && ch <= 'z' {
				ch = ch - 'a' + 'A'
				matrix[i][j] = ch
			}
			switch {
			case i == j:
				if ch != 'X' {
					return fmt.Errorf("diagonal cell (%d,%d) must be 'X'", i+1, j+1)
				}
			case ch == 'W':
				if matrix[j][i] != 'L' && !(matrix[j][i] == 'l') {
					return fmt.Errorf("cells (%d,%d)='W' and (%d,%d) must be 'L'", i+1, j+1, j+1, i+1)
				}
				points[i] += 2
			case ch == 'L':
				if matrix[j][i] != 'W' && !(matrix[j][i] == 'w') {
					return fmt.Errorf("cells (%d,%d)='L' and (%d,%d) must be 'W'", i+1, j+1, j+1, i+1)
				}
			case ch == 'D':
				if matrix[j][i] != 'D' && !(matrix[j][i] == 'd') {
					return fmt.Errorf("cells (%d,%d)='D' require (%d,%d)='D'", i+1, j+1, j+1, i+1)
				}
				points[i]++
			default:
				return fmt.Errorf("invalid character %q at (%d,%d)", ch, i+1, j+1)
			}
		}
	}
	sorted := append([]int(nil), points...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] > sorted[j] })
	for i := 0; i < tc.n; i++ {
		if sorted[i] != tc.top[i] {
			return fmt.Errorf("top score mismatch at position %d: expected %d got %d", i+1, tc.top[i], sorted[i])
		}
	}
	return nil
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	tests := []testCase{
		{m: 1, n: 1, top: []int{0}},
		{m: 2, n: 1, top: []int{2}},
		{m: 2, n: 2, top: []int{2, 0}},
		{m: 3, n: 2, top: []int{4, 2}},
	}
	for i := 0; i < 60; i++ {
		tests = append(tests, randomValidCase(rng, 8))
	}
	for i := 0; i < 60; i++ {
		tests = append(tests, randomGeneralCase(rng, 10))
	}
	for i := 0; i < 30; i++ {
		tests = append(tests, randomGeneralCase(rng, 20))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	if bin == "--" {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := genTests()
	for i, tc := range tests {
		in := inputString(tc)
		refOut, err := runProgram(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
			os.Exit(1)
		}
		refVerdict, refRows, err := parseOutput(refOut, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference parse error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
			os.Exit(1)
		}
		expectedYes := refVerdict == "yes"
		if expectedYes {
			if err := checkTable(tc, refRows); err != nil {
				fmt.Fprintf(os.Stderr, "reference produced invalid table on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
				os.Exit(1)
			}
		}

		out, runErr := runProgram(bin, in)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\ninput:\n%soutput:\n%s\n", i+1, runErr, in, out)
			os.Exit(1)
		}
		verdict, rows, err := parseOutput(out, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
		if verdict == "no" {
			if expectedYes {
				fmt.Fprintf(os.Stderr, "test %d failed: expected 'yes' but got 'no'\ninput:\n%soutput:\n%s\n", i+1, in, out)
				os.Exit(1)
			}
			continue
		}
		if verdict != "yes" {
			fmt.Fprintf(os.Stderr, "test %d invalid verdict %q\ninput:\n%soutput:\n%s\n", i+1, verdict, in, out)
			os.Exit(1)
		}
		if !expectedYes {
			fmt.Fprintf(os.Stderr, "test %d failed: expected 'no' but got 'yes'\ninput:\n%soutput:\n%s\n", i+1, in, out)
			os.Exit(1)
		}
		if err := checkTable(tc, rows); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

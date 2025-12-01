package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource2057G = "2000-2999/2000-2099/2050-2059/2057/2057G.go"

type testCase struct {
	n, m int
	grid []string
}

type namedCase struct {
	name string
	tc   testCase
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource2057G)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	for idx, tc := range tests {
		input := buildInput(tc.tc)

		refOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, input)
			os.Exit(1)
		}
		if err := validateOutputs(tc.tc, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "internal error: reference produced invalid output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinary(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, input)
			os.Exit(1)
		}
		if err := validateOutputs(tc.tc, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier2057G-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(path))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildTests() []namedCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return []namedCase{
		{name: "sample1", tc: testCase{n: 3, m: 3, grid: []string{".#.", "###", ".#."}}},
		{name: "sample2", tc: testCase{n: 2, m: 6, grid: []string{"######", "######"}}},
		{name: "sample3", tc: testCase{n: 3, m: 7, grid: []string{"###....#", ".#.####", "##....#"}}},
		{name: "all_blocked", tc: testCase{n: 2, m: 2, grid: []string{"..", ".."}}},
		{name: "all_free", tc: testCase{n: 4, m: 4, grid: []string{"####", "####", "####", "####"}}},
		{name: "checker", tc: generateChecker(6, 7)},
		{name: "random_small", tc: randomGrid(rng, 10, 12)},
		{name: "random_medium", tc: randomGrid(rng, 100, 120)},
		{name: "large_sparse", tc: randomSparse(rng, 600, 700, 0.1)},
		{name: "large_dense", tc: randomSparse(rng, 700, 600, 0.9)},
		{name: "near_limit", tc: randomSparse(rng, 1000, 1000, 0.5)}, // 1e6 cells
	}
}

func randomGrid(rng *rand.Rand, n, m int) testCase {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(3) == 0 {
				row[j] = '.'
			} else {
				row[j] = '#'
			}
		}
		grid[i] = string(row)
	}
	return testCase{n: n, m: m, grid: grid}
}

func randomSparse(rng *rand.Rand, n, m int, prob float64) testCase {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Float64() < prob {
				row[j] = '#'
			} else {
				row[j] = '.'
			}
		}
		grid[i] = string(row)
	}
	return testCase{n: n, m: m, grid: grid}
}

func generateChecker(n, m int) testCase {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if (i+j)%2 == 0 {
				row[j] = '#'
			} else {
				row[j] = '.'
			}
		}
		grid[i] = string(row)
	}
	return testCase{n: n, m: m, grid: grid}
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, 1)
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, row := range tc.grid {
		sb.WriteString(row)
		if i+1 != len(tc.grid) {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

func validateOutputs(tc testCase, output string) error {
	lines := splitNonEmptyLines(output)
	if len(lines) != tc.n {
		return fmt.Errorf("expected %d lines, got %d", tc.n, len(lines))
	}
	for i, line := range lines {
		if len(line) != tc.m {
			return fmt.Errorf("line %d length mismatch: expected %d got %d", i+1, tc.m, len(line))
		}
	}

	s, p := computeSP(tc)
	maxS := s + p
	maxS = maxS / 5
	if maxS*5 < s+p {
		// floor already used, inequality is cnt*5 <= s+p
	}

	countS := 0
	free := make([][]bool, tc.n)
	sCell := make([][]bool, tc.n)
	for i := 0; i < tc.n; i++ {
		free[i] = make([]bool, tc.m)
		sCell[i] = make([]bool, tc.m)
		for j := 0; j < tc.m; j++ {
			orig := tc.grid[i][j]
			outc := lines[i][j]
			if orig == '.' && outc != '.' {
				return fmt.Errorf("blocked cell (%d,%d) changed to %c", i+1, j+1, outc)
			}
			if orig == '#' && outc == '.' {
				return fmt.Errorf("free cell (%d,%d) incorrectly output as '.'", i+1, j+1)
			}
			if orig == '#' {
				free[i][j] = true
				if outc == 'S' {
					sCell[i][j] = true
					countS++
				} else if outc != '#' {
					return fmt.Errorf("free cell (%d,%d) has invalid char %c", i+1, j+1, outc)
				}
			} else if outc != '.' {
				return fmt.Errorf("blocked cell (%d,%d) has invalid char %c", i+1, j+1, outc)
			}
		}
	}

	if countS*5 > s+p {
		return fmt.Errorf("|S|=%d exceeds limit floor((s+p)/5)=%d", countS, maxS)
	}

	// coverage check
	dir := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if !free[i][j] {
				continue
			}
			if sCell[i][j] {
				continue
			}
			ok := false
			for _, d := range dir {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < tc.n && nj >= 0 && nj < tc.m && sCell[ni][nj] {
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("free cell (%d,%d) not covered by S", i+1, j+1)
			}
		}
	}

	return nil
}

func computeSP(tc testCase) (int, int) {
	s := 0
	p := 0
	n, m := tc.n, tc.m
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if tc.grid[i][j] != '#' {
				continue
			}
			s++
			if i == 0 || tc.grid[i-1][j] != '#' {
				p++
			}
			if i == n-1 || tc.grid[i+1][j] != '#' {
				p++
			}
			if j == 0 || tc.grid[i][j-1] != '#' {
				p++
			}
			if j == m-1 || tc.grid[i][j+1] != '#' {
				p++
			}
		}
	}
	return s, p
}

func splitNonEmptyLines(out string) []string {
	lines := strings.Split(out, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

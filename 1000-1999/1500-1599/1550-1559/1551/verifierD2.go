package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type caseSpec struct {
	n, m, k int
}

type caseResult struct {
	possible bool
	grid     []string
}

type testCase struct {
	desc  string
	specs []caseSpec
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Printf("Reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, refOut)
			os.Exit(1)
		}
		refRes, err := parseOutput(refOut, tc.specs)
		if err != nil {
			fmt.Printf("Failed to parse reference output on test %d (%s): %v\nOutput:\n%s\n", i+1, tc.desc, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		candRes, err := parseOutput(out, tc.specs)
		if err != nil {
			fmt.Printf("Failed to parse output on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		for idx, spec := range tc.specs {
			if !refRes[idx].possible {
				if candRes[idx].possible {
					fmt.Printf("Wrong answer on test %d (%s) case %d: should be NO\nInput:\n%s", i+1, tc.desc, idx+1, tc.input)
					os.Exit(1)
				}
				continue
			}
			if !candRes[idx].possible {
				fmt.Printf("Wrong answer on test %d (%s) case %d: solution exists but NO reported\nInput:\n%s", i+1, tc.desc, idx+1, tc.input)
				os.Exit(1)
			}
			if err := validateGrid(spec, candRes[idx].grid); err != nil {
				fmt.Printf("Invalid layout on test %d (%s) case %d: %v\nInput:\n%s", i+1, tc.desc, idx+1, err, tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	path := "./ref1551D2.bin"
	cmd := exec.Command("go", "build", "-o", path, "1551D2.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	return path, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func parseOutput(out string, specs []caseSpec) ([]caseResult, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	results := make([]caseResult, len(specs))
	for i, spec := range specs {
		var verdict string
		if _, err := fmt.Fscan(reader, &verdict); err != nil {
			return nil, fmt.Errorf("case %d: failed to read verdict: %v", i+1, err)
		}
		v := strings.ToUpper(verdict)
		switch v {
		case "NO":
			results[i] = caseResult{possible: false}
		case "YES":
			rows := make([]string, spec.n)
			for r := 0; r < spec.n; r++ {
				if _, err := fmt.Fscan(reader, &rows[r]); err != nil {
					return nil, fmt.Errorf("case %d: failed to read row %d: %v", i+1, r+1, err)
				}
				if len(rows[r]) != spec.m {
					return nil, fmt.Errorf("case %d: row %d has length %d expected %d", i+1, r+1, len(rows[r]), spec.m)
				}
			}
			results[i] = caseResult{possible: true, grid: rows}
		default:
			return nil, fmt.Errorf("case %d: invalid verdict %q", i+1, verdict)
		}
	}
	if rest := strings.TrimSpace(readRemaining(reader)); rest != "" {
		return nil, fmt.Errorf("unexpected extra output: %q", rest)
	}
	return results, nil
}

func validateGrid(spec caseSpec, grid []string) error {
	if len(grid) != spec.n {
		return fmt.Errorf("expected %d rows got %d", spec.n, len(grid))
	}
	n, m := spec.n, spec.m
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	horizontal := 0
	for i := 0; i < n; i++ {
		if len(grid[i]) != m {
			return fmt.Errorf("row %d length mismatch", i+1)
		}
		for j := 0; j < m; j++ {
			if visited[i][j] {
				continue
			}
			ch := grid[i][j]
			pi, pj := -1, -1
			if j+1 < m && !visited[i][j+1] && grid[i][j+1] == ch {
				pi, pj = i, j+1
				horizontal++
			} else if i+1 < n && !visited[i+1][j] && grid[i+1][j] == ch {
				pi, pj = i+1, j
			} else {
				return fmt.Errorf("cell (%d,%d) does not form a domino", i+1, j+1)
			}
			if err := ensureNoExtraNeighbor(grid, i, j, pi, pj, ch); err != nil {
				return err
			}
			if err := ensureNoExtraNeighbor(grid, pi, pj, i, j, ch); err != nil {
				return err
			}
			visited[i][j] = true
			visited[pi][pj] = true
		}
	}
	if horizontal != spec.k {
		return fmt.Errorf("expected %d horizontal dominos got %d", spec.k, horizontal)
	}
	return nil
}

func ensureNoExtraNeighbor(grid []string, i, j, pi, pj int, ch byte) error {
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	n := len(grid)
	m := len(grid[0])
	for _, d := range dirs {
		ni, nj := i+d[0], j+d[1]
		if ni < 0 || ni >= n || nj < 0 || nj >= m {
			continue
		}
		if ni == pi && nj == pj {
			continue
		}
		if grid[ni][nj] == grid[i][j] {
			return fmt.Errorf("cells (%d,%d) and (%d,%d) illegally share letter %c", i+1, j+1, ni+1, nj+1, ch)
		}
	}
	return nil
}

func readRemaining(r *bufio.Reader) string {
	var sb strings.Builder
	for {
		line, err := r.ReadString('\n')
		sb.WriteString(line)
		if err != nil {
			break
		}
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	add := func(desc string, specs []caseSpec) {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(specs))
		for _, s := range specs {
			fmt.Fprintf(&sb, "%d %d %d\n", s.n, s.m, s.k)
		}
		tests = append(tests, testCase{
			desc:  desc,
			specs: specs,
			input: sb.String(),
		})
	}

	add("basic-small", []caseSpec{{1, 2, 1}, {2, 2, 0}, {2, 3, 1}})
	add("odd-row", []caseSpec{{3, 4, 4}, {5, 2, 4}})
	add("odd-column", []caseSpec{{4, 3, 4}, {6, 5, 8}})

	rng := rand.New(rand.NewSource(123456789))
	for len(tests) < 40 {
		numCases := rng.Intn(5) + 1
		specs := make([]caseSpec, numCases)
		for i := 0; i < numCases; i++ {
			n := rng.Intn(8) + 1
			m := rng.Intn(8) + 1
			if (n*m)%2 == 1 {
				if n < 8 {
					n++
				} else {
					m++
				}
			}
			k := rng.Intn(n*m/2 + 1)
			specs[i] = caseSpec{n: n, m: m, k: k}
		}
		add(fmt.Sprintf("random-small-%d", len(tests)), specs)
	}

	// Large stress tests
	add("large-even", []caseSpec{{n: 100, m: 100, k: 2500}, {n: 100, m: 50, k: 1200}})
	add("large-odd-row", []caseSpec{{n: 99, m: 100, k: 3300}})
	add("large-odd-column", []caseSpec{{n: 100, m: 99, k: 4000}})

	return tests
}

package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const testCount = 120

type testCase struct {
	n, m int
	k    int
}

type parsedResult struct {
	impossible bool
	grid       []string
}

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "720C.go")
	tmp, err := os.CreateTemp("", "oracle720C")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()
	os.Remove(path)
	cmd := exec.Command("go", "build", "-o", path, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	return path, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimRight(out.String(), "\n"), nil
}

func genInput(r *rand.Rand) (string, []testCase) {
	t := r.Intn(4) + 1
	cases := make([]testCase, 0, t)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := 3 + r.Intn(98) // up to 100
		m := 3 + r.Intn(98)
		maxK := 4 * (n - 1) * (m - 1)
		k := 0
		if r.Intn(5) == 0 {
			k = maxK + r.Intn(50) + 1 // force impossible sometimes
		} else if maxK > 0 {
			k = r.Intn(maxK + 1)
		}
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
		cases = append(cases, testCase{n: n, m: m, k: k})
	}
	return sb.String(), cases
}

func stripCR(s string) string {
	return strings.TrimRight(s, "\r")
}

func parseOutput(out string, tests []testCase) ([]parsedResult, error) {
	lines := strings.Split(out, "\n")
	idx := 0
	res := make([]parsedResult, len(tests))
	for tIdx, tc := range tests {
		for idx < len(lines) && strings.TrimSpace(stripCR(lines[idx])) == "" {
			idx++
		}
		if idx >= len(lines) {
			return nil, fmt.Errorf("missing output for test %d", tIdx+1)
		}
		line := stripCR(lines[idx])
		idx++
		if strings.TrimSpace(line) == "-1" {
			res[tIdx] = parsedResult{impossible: true}
			continue
		}
		grid := make([]string, tc.n)
		if len(line) != tc.m {
			return nil, fmt.Errorf("test %d row 1 length %d expected %d", tIdx+1, len(line), tc.m)
		}
		grid[0] = line
		for row := 1; row < tc.n; row++ {
			if idx >= len(lines) {
				return nil, fmt.Errorf("test %d missing row %d", tIdx+1, row+1)
			}
			rowStr := stripCR(lines[idx])
			idx++
			if len(rowStr) != tc.m {
				return nil, fmt.Errorf("test %d row %d length %d expected %d", tIdx+1, row+1, len(rowStr), tc.m)
			}
			grid[row] = rowStr
		}
		res[tIdx] = parsedResult{grid: grid}
	}
	for idx < len(lines) {
		if strings.TrimSpace(stripCR(lines[idx])) != "" {
			return nil, errors.New("extra output data")
		}
		idx++
	}
	return res, nil
}

func validateGrid(tc testCase, grid []string) error {
	if len(grid) != tc.n {
		return fmt.Errorf("expected %d rows, got %d", tc.n, len(grid))
	}
	stars := 0
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	visited := make([][]bool, tc.n)
	for i := range visited {
		visited[i] = make([]bool, tc.m)
	}
	sx, sy := -1, -1
	for i := 0; i < tc.n; i++ {
		if len(grid[i]) != tc.m {
			return fmt.Errorf("row %d length mismatch", i+1)
		}
		for j := 0; j < tc.m; j++ {
			c := grid[i][j]
			if c != '.' && c != '*' {
				return fmt.Errorf("invalid char %q at (%d,%d)", c, i+1, j+1)
			}
			if c == '*' {
				stars++
				if sx == -1 {
					sx, sy = i, j
				}
			}
		}
	}
	if stars == 0 {
		return errors.New("no marked squares")
	}
	// BFS for connectivity
	queue := make([][2]int, 0, stars)
	queue = append(queue, [2]int{sx, sy})
	visited[sx][sy] = true
	visitedStars := 1
	for idx := 0; idx < len(queue); idx++ {
		x, y := queue[idx][0], queue[idx][1]
		for _, d := range dirs {
			nx, ny := x+d[0], y+d[1]
			if nx < 0 || ny < 0 || nx >= tc.n || ny >= tc.m {
				continue
			}
			if grid[nx][ny] != '*' || visited[nx][ny] {
				continue
			}
			visited[nx][ny] = true
			visitedStars++
			queue = append(queue, [2]int{nx, ny})
		}
	}
	if visitedStars != stars {
		return errors.New("marked squares not connected")
	}
	// count L-trominoes
	var total int64
	for i := 0; i < tc.n-1; i++ {
		for j := 0; j < tc.m-1; j++ {
			cnt := 0
			if grid[i][j] == '*' {
				cnt++
			}
			if grid[i+1][j] == '*' {
				cnt++
			}
			if grid[i][j+1] == '*' {
				cnt++
			}
			if grid[i+1][j+1] == '*' {
				cnt++
			}
			if cnt == 3 {
				total++
			} else if cnt == 4 {
				total += 4
			}
		}
	}
	if total != int64(tc.k) {
		return fmt.Errorf("expected %d L-trominoes, got %d", tc.k, total)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	r := rand.New(rand.NewSource(1))
	for t := 0; t < testCount; t++ {
		input, cases := genInput(r)
		expectStr, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		expectParsed, err := parseOutput(expectStr, cases)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotParsed, err := parseOutput(gotStr, cases)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		for i, tc := range cases {
			if expectParsed[i].impossible {
				if gotParsed[i].impossible {
					continue
				}
				fmt.Printf("test %d case %d failed\ninput:\n%s\nerror: expected -1 but found grid\n", t+1, i+1, input)
				os.Exit(1)
			}
			if gotParsed[i].impossible {
				fmt.Printf("test %d case %d failed\ninput:\n%s\nerror: unexpected -1\n", t+1, i+1, input)
				os.Exit(1)
			}
			if err := validateGrid(tc, gotParsed[i].grid); err != nil {
				fmt.Printf("test %d case %d failed\ninput:\n%s\nerror: %v\n", t+1, i+1, input, err)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}

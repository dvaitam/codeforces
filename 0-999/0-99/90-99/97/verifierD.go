package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	input   string
	expect  string
	comment string
}

type cell struct {
	x, y int
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]
	if candidate == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	baseDir := currentDir()
	refBin, err := buildReference(baseDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	inputs := generateInputs(rng)
	tests := make([]testCase, len(inputs))
	for i, inp := range inputs {
		exp, err := runProgram(refBin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d: %v\ninput:\n%s", i+1, err, inp)
			os.Exit(1)
		}
		tests[i] = testCase{
			input:  inp,
			expect: exp,
		}
	}

	for i, tc := range tests {
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(tc.expect) {
			fmt.Fprintf(os.Stderr, "wrong answer on case %d\nInput:\n%sExpected: %s\nGot: %s\n", i+1, tc.input, tc.expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func currentDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine caller info")
	}
	return filepath.Dir(file)
}

func buildReference(baseDir string) (string, error) {
	refBin := filepath.Join(baseDir, "ref97D.bin")
	cmd := exec.Command("go", "build", "-o", refBin, "97D.go")
	cmd.Dir = baseDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	return refBin, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateInputs(rng *rand.Rand) []string {
	var tests []string
	tests = append(tests, singleExitCase())
	tests = append(tests, impossibleTwoCellsCase())
	tests = append(tests, corridorCase())
	for len(tests) < 40 {
		tests = append(tests, randomCase(rng, 3, 12, 3, 12, 200))
	}
	for i := 0; i < 20; i++ {
		tests = append(tests, randomCase(rng, 5, 25, 5, 25, 2000))
	}
	for i := 0; i < 3; i++ {
		tests = append(tests, randomCaseWithFixedK(rng, 10, 30, 10, 30, 100000))
	}
	return tests
}

func singleExitCase() string {
	grid := []string{
		"###",
		"#E#",
		"###",
	}
	cmds := strings.Repeat("L", 5)
	return formatInput(3, 3, cmds, grid)
}

func impossibleTwoCellsCase() string {
	grid := []string{
		"####",
		"#E.#",
		"####",
		"####",
	}
	cmds := strings.Repeat("R", 6)
	return formatInput(4, 4, cmds, grid)
}

func corridorCase() string {
	grid := []string{
		"#####",
		"#E..#",
		"#####",
		"#####",
		"#####",
	}
	cmds := "RRLL"
	return formatInput(5, 5, cmds, grid)
}

func randomCase(rng *rand.Rand, minN, maxN, minM, maxM, maxK int) string {
	n := rng.Intn(maxN-minN+1) + minN
	m := rng.Intn(maxM-minM+1) + minM
	k := rng.Intn(maxK) + 1
	grid := randomGrid(rng, n, m, 0.35)
	cmds := randomCommands(rng, k)
	return formatInput(n, m, cmds, grid)
}

func randomCaseWithFixedK(rng *rand.Rand, minN, maxN, minM, maxM, k int) string {
	n := rng.Intn(maxN-minN+1) + minN
	m := rng.Intn(maxM-minM+1) + minM
	grid := randomGrid(rng, n, m, 0.25)
	cmds := randomCommands(rng, k)
	return formatInput(n, m, cmds, grid)
}

func randomCommands(rng *rand.Rand, k int) string {
	dir := []byte{'L', 'R', 'U', 'D'}
	var sb strings.Builder
	for i := 0; i < k; i++ {
		sb.WriteByte(dir[rng.Intn(len(dir))])
	}
	return sb.String()
}

func formatInput(n, m int, cmds string, grid []string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, len(cmds))
	for _, row := range grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	sb.WriteString(cmds)
	sb.WriteByte('\n')
	return sb.String()
}

func randomGrid(rng *rand.Rand, n, m int, wallProb float64) []string {
	for attempts := 0; attempts < 200; attempts++ {
		grid := make([][]byte, n)
		for i := range grid {
			grid[i] = make([]byte, m)
			for j := range grid[i] {
				grid[i][j] = '#'
			}
		}
		var interior []cell
		for i := 1; i < n-1; i++ {
			for j := 1; j < m-1; j++ {
				grid[i][j] = '.'
				interior = append(interior, cell{i, j})
			}
		}
		if len(interior) == 0 {
			continue
		}
		exitIdx := rng.Intn(len(interior))
		ex := interior[exitIdx]
		grid[ex.x][ex.y] = 'E'
		for idx, c := range interior {
			if idx == exitIdx {
				continue
			}
			if rng.Float64() < wallProb {
				grid[c.x][c.y] = '#'
			}
		}
		if ensuredConnected(grid) {
			rows := make([]string, n)
			for i := 0; i < n; i++ {
				rows[i] = string(grid[i])
			}
			return rows
		}
	}
	// fallback to simple grid with all interior passable
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if i == 0 || j == 0 || i == n-1 || j == m-1 {
				row[j] = '#'
			} else if i == 1 && j == 1 {
				row[j] = 'E'
			} else {
				row[j] = '.'
			}
		}
		rows[i] = string(row)
	}
	return rows
}

func ensuredConnected(grid [][]byte) bool {
	n := len(grid)
	m := len(grid[0])
	total := 0
	var sx, sy int
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] != '#' {
				total++
				if grid[i][j] == 'E' {
					sx, sy = i, j
				}
			}
		}
	}
	if total == 0 {
		return false
	}
	visited := make([][]bool, n)
	for i := 0; i < n; i++ {
		visited[i] = make([]bool, m)
	}
	queue := []cell{{sx, sy}}
	visited[sx][sy] = true
	seen := 0
	dirs := []cell{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		seen++
		for _, d := range dirs {
			nx := cur.x + d.x
			ny := cur.y + d.y
			if nx < 0 || ny < 0 || nx >= n || ny >= m {
				continue
			}
			if grid[nx][ny] == '#' || visited[nx][ny] {
				continue
			}
			visited[nx][ny] = true
			queue = append(queue, cell{nx, ny})
		}
	}
	return seen == total
}

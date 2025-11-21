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

type testCase struct {
	name  string
	input string
}

type point struct{ x, y int }

func solveRef(input string) ([]int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return nil, err
	}
	grid := make([][]bool, n)
	for i := 0; i < n; i++ {
		line := ""
		for len(line) < m {
			part, err := reader.ReadString('\n')
			if err != nil && len(part) == 0 {
				break
			}
			line += part
		}
		row := make([]bool, m)
		for j := 0; j < m; j++ {
			if line[j] == '.' {
				row[j] = true
			}
		}
		grid[i] = row
	}
	xs := make([]int, k)
	ys := make([]int, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &xs[i], &ys[i])
		xs[i]--
		ys[i]--
		grid[xs[i]][ys[i]] = false
	}
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, m)
	}
	maxSq := 0
	min3 := func(a, b, c int) int {
		if a < b {
			if a < c {
				return a
			}
			return c
		}
		if b < c {
			return b
		}
		return c
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] {
				if i > 0 && j > 0 {
					dp[i][j] = 1 + min3(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
				} else {
					dp[i][j] = 1
				}
				if dp[i][j] > maxSq {
					maxSq = dp[i][j]
				}
			}
		}
	}
	ans := make([]int, k+1)
	ans[k] = maxSq
	queue := make([]point, 0, n*m)
	for t := k - 1; t >= 0; t-- {
		i := xs[t]
		j := ys[t]
		grid[i][j] = true
		newVal := 1
		if i > 0 && j > 0 {
			newVal = 1 + min3(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
		}
		if newVal > dp[i][j] {
			dp[i][j] = newVal
			if newVal > maxSq {
				maxSq = newVal
			}
			queue = queue[:0]
			queue = append(queue, point{i, j})
			for head := 0; head < len(queue); head++ {
				cur := queue[head]
				dirs := [][2]int{{1, 0}, {0, 1}, {1, 1}}
				for _, d := range dirs {
					ni := cur.x + d[0]
					nj := cur.y + d[1]
					if ni < n && nj < m && grid[ni][nj] {
						cand := 1
						if ni > 0 && nj > 0 {
							cand = 1 + min3(dp[ni-1][nj], dp[ni][nj-1], dp[ni-1][nj-1])
						}
						if cand > dp[ni][nj] {
							dp[ni][nj] = cand
							if cand > maxSq {
								maxSq = cand
							}
							queue = append(queue, point{ni, nj})
						}
					}
				}
			}
		}
		ans[t] = maxSq
	}
	return ans[1:], nil
}

func makeCase(name string, n, m int, grid []string, ops []point) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, len(ops)))
	for _, row := range grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	for _, op := range ops {
		sb.WriteString(fmt.Sprintf("%d %d\n", op.x+1, op.y+1))
	}
	return testCase{name: name, input: sb.String()}
}

func randomGrid(rng *rand.Rand, n, m int) []string {
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(4) == 0 {
				row[j] = 'X'
			} else {
				row[j] = '.'
			}
		}
		grid[i] = string(row)
	}
	return grid
}

func randomOps(rng *rand.Rand, grid [][]byte, k int) ([]point, []string) {
	n := len(grid)
	m := len(grid[0])
	ops := make([]point, 0, k)
	for len(ops) < k {
		x := rng.Intn(n)
		y := rng.Intn(m)
		if grid[x][y] == '.' {
			grid[x][y] = 'X'
			ops = append(ops, point{x, y})
		}
	}
	strGrid := make([]string, n)
	for i := 0; i < n; i++ {
		strGrid[i] = string(grid[i])
	}
	return ops, strGrid
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(480))
	var tests []testCase
	gen := func(prefix string, count, maxN, maxM, maxK int) {
		for i := 0; i < count; i++ {
			n := rng.Intn(maxN) + 1
			m := rng.Intn(maxM) + 1
			grid := make([][]byte, n)
			for r := 0; r < n; r++ {
				grid[r] = []byte(randomGrid(rng, 1, m)[0])
			}
			k := rng.Intn(maxK) + 1
			ops, finalGrid := randomOps(rng, grid, k)
			tests = append(tests, makeCase(fmt.Sprintf("%s_%d", prefix, i+1), n, m, finalGrid, ops))
		}
	}
	gen("small", 60, 3, 3, 4)
	gen("medium", 60, 8, 8, 10)
	gen("large", 40, 15, 15, 40)
	return tests
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single_cell", 1, 1, []string{"."}, []point{{0, 0}}),
		makeCase("row", 1, 5, []string{"....."}, []point{{0, 2}, {0, 1}, {0, 3}}),
		makeCase("column", 5, 1, []string{".", ".", ".", ".", "."}, []point{{2, 0}, {4, 0}}),
		makeCase("preblocked", 3, 3, []string{"XXX", "X.X", "XXX"}, []point{{1, 1}}),
		makeCase("checker", 4, 4, []string{
			".X.X",
			"X.X.",
			".X.X",
			"X.X.",
		}, []point{{0, 0}, {1, 1}, {2, 2}, {3, 3}}),
	}
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseOutput(output string, expectedLen int) ([]int, error) {
	fields := strings.Fields(output)
	if len(fields) != expectedLen {
		return nil, fmt.Errorf("expected %d integers, got %d", expectedLen, len(fields))
	}
	ans := make([]int, expectedLen)
	for i, f := range fields {
		if _, err := fmt.Sscan(f, &ans[i]); err != nil {
			return nil, fmt.Errorf("failed to parse integer #%d: %v", i+1, err)
		}
	}
	return ans, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to build reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		output, runErr := runCandidate(bin, tc.input)
		if runErr != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, runErr, tc.input)
			os.Exit(1)
		}
		got, parseErr := parseOutput(output, len(expect))
		if parseErr != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, parseErr, tc.input, output)
			os.Exit(1)
		}
		match := true
		for i := range expect {
			if expect[i] != got[i] {
				match = false
				break
			}
		}
		if !match {
			fmt.Printf("test %d (%s) mismatch\ninput:\n%s\nexpect:%v\nactual:%v\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

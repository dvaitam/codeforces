package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseE struct {
	n, m, k, d int
	grid       [][]int
	exp        string
}

func rowCost(row []int, d int) int64 {
	m := len(row)
	const INF int64 = 1 << 60
	cost := make([]int64, m)
	for i, v := range row {
		cost[i] = int64(v) + 1
	}
	dp := make([]int64, m)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = cost[0]
	type pair struct {
		idx int
		val int64
	}
	deque := make([]pair, 0)
	for j := 1; j < m; j++ {
		val := dp[j-1]
		for len(deque) > 0 && deque[len(deque)-1].val >= val {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, pair{j - 1, val})
		limit := j - (d + 1)
		for len(deque) > 0 && deque[0].idx < limit {
			deque = deque[1:]
		}
		if len(deque) > 0 {
			dp[j] = cost[j] + deque[0].val
		}
	}
	return dp[m-1]
}

func solveE(n, m, k, d int, grid [][]int) string {
	costs := make([]int64, n)
	for i := 0; i < n; i++ {
		costs[i] = rowCost(grid[i], d)
	}
	var sum int64
	for i := 0; i < k; i++ {
		sum += costs[i]
	}
	ans := sum
	for i := k; i < n; i++ {
		sum += costs[i] - costs[i-k]
		if sum < ans {
			ans = sum
		}
	}
	return fmt.Sprint(ans)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func generateTests() []testCaseE {
	rng := rand.New(rand.NewSource(5))
	cases := make([]testCaseE, 100)
	for i := range cases {
		n := rng.Intn(4) + 1
		m := rng.Intn(6) + 3
		k := rng.Intn(n) + 1
		d := rng.Intn(m) + 1
		grid := make([][]int, n)
		for r := 0; r < n; r++ {
			grid[r] = make([]int, m)
			for c := 0; c < m; c++ {
				if c == 0 || c == m-1 {
					grid[r][c] = 0
				} else {
					grid[r][c] = rng.Intn(10)
				}
			}
		}
		cases[i] = testCaseE{n: n, m: m, k: k, d: d, grid: grid, exp: solveE(n, m, k, d, grid)}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintln(&sb, 1)
		fmt.Fprintf(&sb, "%d %d %d %d\n", tc.n, tc.m, tc.k, tc.d)
		for r := 0; r < tc.n; r++ {
			for c := 0; c < tc.m; c++ {
				if c > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprint(tc.grid[r][c]))
			}
			sb.WriteByte('\n')
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

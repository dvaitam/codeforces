package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type obstacle struct{ x, y int }

type testCaseB struct {
	n        int
	obs      []obstacle
	expected string
}

func solveCase(tc testCaseB) string {
	n := tc.n
	grid := make([][]bool, n+1)
	for i := range grid {
		grid[i] = make([]bool, n+1)
	}
	for _, o := range tc.obs {
		if o.x >= 1 && o.x <= n && o.y >= 1 && o.y <= n {
			grid[o.x][o.y] = true
		}
	}
	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, n+1)
	}
	dp[1][1] = !grid[1][1]
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if grid[i][j] {
				dp[i][j] = false
				continue
			}
			if i == 1 && j == 1 {
				continue
			}
			if i > 1 && dp[i-1][j] {
				dp[i][j] = true
			}
			if j > 1 && dp[i][j-1] {
				dp[i][j] = true
			}
		}
	}
	if dp[n][n] {
		return strconv.Itoa(2*n - 2)
	}
	return "-1"
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 2
		maxObs := n * n / 3
		m := rng.Intn(maxObs + 1)
		obsSet := make(map[[2]int]struct{})
		obs := make([]obstacle, 0, m)
		for len(obs) < m {
			x := rng.Intn(n) + 1
			y := rng.Intn(n) + 1
			if x == 1 && y == 1 || x == n && y == n {
				continue
			}
			key := [2]int{x, y}
			if _, ok := obsSet[key]; ok {
				continue
			}
			obsSet[key] = struct{}{}
			obs = append(obs, obstacle{x, y})
		}
		tc := testCaseB{n: n, obs: obs}
		expected := solveCase(tc)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, len(obs)))
		for _, o := range obs {
			sb.WriteString(fmt.Sprintf("%d %d\n", o.x, o.y))
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

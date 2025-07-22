package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve(n, m, x, y int, grid []string) int {
	blacks := make([]int, m+1)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				blacks[j+1]++
			}
		}
	}
	costW := make([]int, m+1)
	costB := make([]int, m+1)
	for j := 1; j <= m; j++ {
		costW[j] = blacks[j]
		costB[j] = n - blacks[j]
	}
	prefW := make([]int, m+1)
	prefB := make([]int, m+1)
	for j := 1; j <= m; j++ {
		prefW[j] = prefW[j-1] + costW[j]
		prefB[j] = prefB[j-1] + costB[j]
	}
	const INF = 1 << 30
	dp := make([][2]int, m+1)
	for j := 0; j <= m; j++ {
		dp[j][0], dp[j][1] = INF, INF
	}
	dp[0][0], dp[0][1] = 0, 0
	for j := 1; j <= m; j++ {
		for k := x; k <= y && k <= j; k++ {
			cost := prefW[j] - prefW[j-k]
			if dp[j-k][1]+cost < dp[j][0] {
				dp[j][0] = dp[j-k][1] + cost
			}
		}
		for k := x; k <= y && k <= j; k++ {
			cost := prefB[j] - prefB[j-k]
			if dp[j-k][0]+cost < dp[j][1] {
				dp[j][1] = dp[j-k][0] + cost
			}
		}
	}
	ans := dp[m][0]
	if dp[m][1] < ans {
		ans = dp[m][1]
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(3))
	for t := 1; t <= 100; t++ {
		n := r.Intn(4) + 3 //3..6
		m := r.Intn(4) + 3 //3..6
		x := r.Intn(min(3, m)) + 1
		y := r.Intn(m-x+1) + x
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			sb := strings.Builder{}
			for j := 0; j < m; j++ {
				if r.Intn(2) == 0 {
					sb.WriteByte('.')
				} else {
					sb.WriteByte('#')
				}
			}
			grid[i] = sb.String()
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d\n", n, m, x, y)
		for i := 0; i < n; i++ {
			sb.WriteString(grid[i])
			sb.WriteByte('\n')
		}
		input := sb.String()
		expected := fmt.Sprintf("%d", solve(n, m, x, y, grid))
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s", t, err, input)
			return
		}
		if strings.TrimSpace(out) != expected {
			fmt.Printf("Test %d FAILED\nInput:\n%sExpected:%s Got:%s\n", t, input, expected, out)
			return
		}
	}
	fmt.Println("All tests passed")
}

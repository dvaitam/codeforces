package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func abs(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func solve(matrix [][]int) int {
	n := len(matrix)
	m := len(matrix[0])
	minimum := make([][]int, n)
	for i := 0; i < n; i++ {
		minimum[i] = make([]int, n)
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			tmp := int(^uint(0) >> 1)
			for k := 0; k < m; k++ {
				d := abs(matrix[i][k], matrix[j][k])
				if d < tmp {
					tmp = d
				}
			}
			minimum[i][j] = tmp
			minimum[j][i] = tmp
		}
	}
	fullMask := (1 << n) - 1
	dp := make([][][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([][]int, n)
		for j := 0; j < n; j++ {
			dp[i][j] = make([]int, 1<<n)
			for k := range dp[i][j] {
				dp[i][j][k] = -1
			}
		}
	}
	var calc func(fi, pre, mask int) int
	calc = func(fi, pre, mask int) int {
		if mask == fullMask {
			r := int(^uint(0) >> 1)
			for i := 0; i < m-1; i++ {
				d := abs(matrix[pre][i], matrix[fi][i+1])
				if d < r {
					r = d
				}
			}
			return r
		}
		if dp[fi][pre][mask] != -1 {
			return dp[fi][pre][mask]
		}
		best := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) == 0 {
				c := calc(fi, i, mask|(1<<i))
				if minimum[pre][i] < c {
					c = minimum[pre][i]
				}
				if c > best {
					best = c
				}
			}
		}
		dp[fi][pre][mask] = best
		return best
	}
	res := 0
	for i := 0; i < n; i++ {
		v := calc(i, i, 1<<i)
		if v > res {
			res = v
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	m := rng.Intn(4) + 2
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, m)
		for j := 0; j < m; j++ {
			matrix[i][j] = rng.Intn(50)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", matrix[i][j]))
		}
		sb.WriteByte('\n')
	}
	ans := solve(matrix)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

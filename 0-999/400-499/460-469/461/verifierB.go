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

const mod = 1000000007

func add(a, b int) int {
	a += b
	if a >= mod {
		a -= mod
	}
	return a
}

func mul(a, b int) int {
	return int((int64(a) * int64(b)) % mod)
}

func solveCase(n int, parents []int, color []int) int {
	children := make([][]int, n)
	for i := 1; i < n; i++ {
		p := parents[i]
		children[p] = append(children[p], i)
	}
	dp := make([][3]int, n)
	for u := 0; u < n; u++ {
		if color[u] == 1 {
			dp[u][0] = 0
			dp[u][1] = 1
			dp[u][2] = 0
		} else {
			dp[u][0] = 1
			dp[u][1] = 0
			dp[u][2] = 0
		}
	}
	for u := n - 1; u >= 0; u-- {
		for _, v := range children[u] {
			t0 := dp[v][0]
			t1 := dp[v][1]
			t2 := dp[v][2]
			cur := dp[u]
			dp[u][0], dp[u][1], dp[u][2] = 0, 0, 0
			for j := 0; j < 3; j++ {
				cj := cur[j]
				if cj == 0 {
					continue
				}
				dp[u][j] = add(dp[u][j], mul(cj, t1))
				if t0 != 0 {
					nj := j
					dp[u][nj] = add(dp[u][nj], mul(cj, t0))
				}
				if t1 != 0 {
					nj := j + 1
					if nj > 2 {
						nj = 2
					}
					dp[u][nj] = add(dp[u][nj], mul(cj, t1))
				}
				if t2 != 0 {
					nj := 2
					dp[u][nj] = add(dp[u][nj], mul(cj, t2))
				}
			}
		}
	}
	return dp[0][1]
}

type caseB struct {
	n       int
	parents []int
	color   []int
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	parents := make([]int, n)
	for i := 1; i < n; i++ {
		parents[i] = rng.Intn(i)
	}
	color := make([]int, n)
	black := 0
	for i := 0; i < n; i++ {
		color[i] = rng.Intn(2)
		if color[i] == 1 {
			black++
		}
	}
	if black == 0 {
		color[rng.Intn(n)] = 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 1; i < n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", parents[i])
	}
	if n > 1 {
		sb.WriteByte('\n')
	} else {
		sb.WriteString("\n")
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", color[i])
	}
	sb.WriteByte('\n')
	ans := solveCase(n, parents, color)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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

package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expected(m int, x int64, c, h []int) int {
	maxH := 0
	for _, v := range h {
		maxH += v
	}
	const inf int64 = math.MaxInt64 / 4
	dp := make([]int64, maxH+1)
	for i := range dp {
		dp[i] = inf
	}
	dp[0] = 0
	for i := 1; i <= m; i++ {
		for j := maxH; j >= 0; j-- {
			if dp[j] == inf {
				continue
			}
			if dp[j]+int64(c[i-1]) <= int64(i-1)*x {
				nj := j + h[i-1]
				if nj <= maxH {
					val := dp[j] + int64(c[i-1])
					if val < dp[nj] {
						dp[nj] = val
					}
				}
			}
		}
	}
	res := 0
	for i := maxH; i >= 0; i-- {
		if dp[i] != inf {
			res = i
			break
		}
	}
	return res
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		m := rng.Intn(5) + 1
		x := int64(rng.Intn(20) + 1)
		c := make([]int, m)
		h := make([]int, m)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", m, x))
		for j := 0; j < m; j++ {
			c[j] = rng.Intn(20)
			h[j] = rng.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", c[j], h[j]))
		}
		input := sb.String()
		exp := fmt.Sprintf("%d", expected(m, x, c, h))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

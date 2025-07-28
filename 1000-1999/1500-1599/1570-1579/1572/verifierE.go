package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func area2(x, y []int64, i, j int) int64 {
	sum := int64(0)
	for t := i; t < j; t++ {
		sum += x[t]*y[t+1] - x[t+1]*y[t]
	}
	sum += x[j]*y[i] - x[i]*y[j]
	if sum < 0 {
		sum = -sum
	}
	return sum
}

func can(x, y []int64, k int, T int64) bool {
	n := len(x)
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	for l := 2; l < n; l++ {
		for i := 0; i+l < n; i++ {
			j := i + l
			if area2(x, y, i, j) >= T {
				dp[i][j] = 1
				for m := i + 1; m < j; m++ {
					if dp[i][m] > 0 && dp[m][j] > 0 {
						if dp[i][m]+dp[m][j] > dp[i][j] {
							dp[i][j] = dp[i][m] + dp[m][j]
						}
					}
				}
			}
		}
	}
	return dp[0][n-1] >= k+1
}

func solveE(reader *bufio.Reader) string {
	var n, k int
	fmt.Fscan(reader, &n, &k)
	x := make([]int64, n)
	y := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &x[i], &y[i])
	}
	total := area2(x, y, 0, n-1)
	low, high := int64(0), total/int64(k+1)
	for low < high {
		mid := (low + high + 1) / 2
		if can(x, y, k, mid) {
			low = mid
		} else {
			high = mid - 1
		}
	}
	return fmt.Sprint(low)
}

type testCaseE struct {
	n int
	k int
	x []int64
	y []int64
}

func generateCaseE(rng *rand.Rand) testCaseE {
	n := rng.Intn(6) + 3 // 3..8
	k := rng.Intn(n - 2)
	x := make([]int64, n)
	y := make([]int64, n)
	for i := 0; i < n; i++ {
		x[i] = int64(rng.Intn(11) - 5)
		y[i] = int64(rng.Intn(11) - 5)
	}
	return testCaseE{n: n, k: k, x: x, y: y}
}

func buildInputE(t testCaseE) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.k))
	for i := 0; i < t.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", t.x[i], t.y[i]))
	}
	return sb.String()
}

func expectedE(t testCaseE) string {
	input := buildInputE(t)
	return solveE(bufio.NewReader(strings.NewReader(input)))
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		input := buildInputE(tc)
		expect := expectedE(tc)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nOutput:%s", i+1, err, out)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		exp := strings.TrimSpace(expect)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

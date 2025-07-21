package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func solveCase(xVals, cVals []int64) int64 {
	n := len(xVals) - 1
	prefixX := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefixX[i] = prefixX[i-1] + xVals[i]
	}
	const INF = int64(9e18)
	dp := make([]int64, n+1)
	dp[1] = cVals[1]
	for i := 2; i <= n; i++ {
		best := INF
		for p := 1; p <= i-1; p++ {
			moveCost := (prefixX[i-1] - prefixX[p]) - int64(i-1-p)*xVals[p]
			if dp[p]+moveCost < best {
				best = dp[p] + moveCost
			}
		}
		dp[i] = best + cVals[i]
	}
	ans := INF
	for i := 1; i <= n; i++ {
		tailCost := (prefixX[n] - prefixX[i]) - int64(n-i)*xVals[i]
		if dp[i]+tailCost < ans {
			ans = dp[i] + tailCost
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	arr := make([]struct{ x, c int64 }, n)
	cur := int64(rng.Intn(10))
	for i := 0; i < n; i++ {
		cur += int64(rng.Intn(10) + 1)
		arr[i].x = cur
		arr[i].c = int64(rng.Intn(20) + 1)
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].x < arr[j].x })
	x := make([]int64, n+1)
	c := make([]int64, n+1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		x[i] = arr[i-1].x
		c[i] = arr[i-1].c
		sb.WriteString(fmt.Sprintf("%d %d\n", x[i], c[i]))
	}
	ans := solveCase(x, c)
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
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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

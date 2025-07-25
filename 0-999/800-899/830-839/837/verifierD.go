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

func expected(n, k int, arr []int64) string {
	two := make([]int, n)
	five := make([]int, n)
	maxFive := 0
	for i := 0; i < n; i++ {
		x := arr[i]
		c2, c5 := 0, 0
		for x%2 == 0 {
			c2++
			x /= 2
		}
		for x%5 == 0 {
			c5++
			x /= 5
		}
		two[i] = c2
		five[i] = c5
		maxFive += c5
	}
	dp := make([][]int, k+1)
	for i := 0; i <= k; i++ {
		dp[i] = make([]int, maxFive+1)
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}
	dp[0][0] = 0
	for idx := 0; idx < n; idx++ {
		t2, f5 := two[idx], five[idx]
		for j := k; j >= 1; j-- {
			for f := maxFive; f >= f5; f-- {
				if dp[j-1][f-f5] >= 0 {
					val := dp[j-1][f-f5] + t2
					if val > dp[j][f] {
						dp[j][f] = val
					}
				}
			}
		}
	}
	ans := 0
	for f := 0; f <= maxFive; f++ {
		val := dp[k][f]
		if val < 0 {
			continue
		}
		if val < f {
			if val > ans {
				ans = val
			}
		} else {
			if f > ans {
				ans = f
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	k := rng.Intn(n) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Int63n(1000000) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", arr[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), expected(n, k, arr)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

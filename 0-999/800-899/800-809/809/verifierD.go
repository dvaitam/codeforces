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

func expectedD(intervals [][2]int64) int {
	const INF int64 = 1 << 60
	n := len(intervals)
	dp := make([]int64, n+2)
	for i := 1; i < len(dp); i++ {
		dp[i] = INF
	}
	dp[0] = -INF
	maxLen := 0
	for _, it := range intervals {
		l, r := it[0], it[1]
		for j := maxLen; j >= 0; j-- {
			if dp[j] == INF {
				continue
			}
			x := dp[j] + 1
			if x < l {
				x = l
			}
			if x <= r && x < dp[j+1] {
				dp[j+1] = x
				if j+1 > maxLen {
					maxLen = j + 1
				}
			}
		}
	}
	return maxLen
}

func generateCaseD(rng *rand.Rand) [][2]int64 {
	n := rng.Intn(10) + 1
	res := make([][2]int64, n)
	base := int64(rng.Intn(50))
	for i := 0; i < n; i++ {
		l := base + int64(rng.Intn(20))
		r := l + int64(rng.Intn(10))
		res[i] = [2]int64{l, r}
	}
	return res
}

func runCase(bin string, intervals [][2]int64) error {
	var input strings.Builder
	n := len(intervals)
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		input.WriteString(fmt.Sprintf("%d %d\n", intervals[i][0], intervals[i][1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(bytes.NewReader(out.Bytes()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	expect := expectedD(intervals)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		intervals := generateCaseD(rng)
		if err := runCase(bin, intervals); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%v\n", i+1, err, intervals)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
